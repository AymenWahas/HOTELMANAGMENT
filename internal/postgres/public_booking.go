package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"hotelmanagement/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) PublicRooms(ctx context.Context, checkIn, checkOut string) ([]domain.PublicRoom, error) {
	query := "SELECT room.id::text, room.room_type, room.rate::float8 " +
		"FROM rooms room WHERE room.archived_at IS NULL AND room.status='available' " +
		"AND $1::date >= current_date AND $2::date > $1::date " +
		"AND NOT EXISTS (SELECT 1 FROM reservations reservation WHERE reservation.room_id=room.id " +
		"AND reservation.archived_at IS NULL AND reservation.status IN ('confirmed','checked_in') " +
		"AND daterange(reservation.check_in,reservation.check_out,'[)') && daterange($1::date,$2::date,'[)')) " +
		"ORDER BY room.rate,room.number"
	rows, err := r.db.Query(ctx, query, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rooms := make([]domain.PublicRoom, 0)
	for rows.Next() {
		var room domain.PublicRoom
		if err := rows.Scan(&room.ID, &room.RoomType, &room.Rate); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *Repository) CreateBookingRequest(ctx context.Context, request domain.BookingRequest) (domain.Record, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Record{}, err
	}
	defer tx.Rollback(ctx)

	var rate float64
	err = tx.QueryRow(ctx, "SELECT rate::float8 FROM rooms WHERE id=$1 AND archived_at IS NULL AND status='available' AND $2::date >= current_date AND $3::date > $2::date FOR UPDATE", request.RoomID, request.CheckIn, request.CheckOut).Scan(&rate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Record{}, ErrRoomUnavailable
		}
		return domain.Record{}, err
	}
	var conflict bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM reservations WHERE room_id=$1 AND archived_at IS NULL AND status IN ('confirmed','checked_in') AND daterange(check_in,check_out,'[)') && daterange($2::date,$3::date,'[)'))",
		request.RoomID, request.CheckIn, request.CheckOut).Scan(&conflict); err != nil {
		return domain.Record{}, err
	}
	if conflict {
		return domain.Record{}, ErrConflict
	}

	var guestID string
	err = tx.QueryRow(ctx, "INSERT INTO guests(first_name,last_name,email,phone) VALUES ($1,$2,$3,$4) "+
		"ON CONFLICT (lower(email)) WHERE email IS NOT NULL AND archived_at IS NULL DO NOTHING RETURNING id::text",
		request.FirstName, request.LastName, request.Email, request.Phone).Scan(&guestID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, "SELECT id::text FROM guests WHERE lower(email)=lower($1) AND archived_at IS NULL", request.Email).Scan(&guestID)
	}
	if err != nil {
		return domain.Record{}, translateConstraint(err)
	}

	item, err := scanRecord(tx.QueryRow(ctx, "INSERT INTO reservations(room_id,guest_id,check_in,check_out,status,adults,children,rate,source) "+
		"VALUES ($1,$2,$3::date,$4::date,'pending',$5,$6,$7,'website') "+
		"RETURNING id::text,'reservation'::text,to_jsonb(reservations)-ARRAY['id','created_at','updated_at','archived_at'],created_at,updated_at",
		request.RoomID, guestID, request.CheckIn, request.CheckOut, request.Adults, request.Children, rate))
	if err != nil {
		return domain.Record{}, translateConstraint(err)
	}
	details, err := json.Marshal(map[string]string{"source": "website", "status": "pending", "email": request.Email})
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := tx.Exec(ctx, "INSERT INTO audit_events(record_kind,record_id,action,details) VALUES ($1,$2,'requested',$3::jsonb)",
		"reservation", item.ID, details); err != nil {
		return domain.Record{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Record{}, translateConstraint(err)
	}
	return item, nil
}
