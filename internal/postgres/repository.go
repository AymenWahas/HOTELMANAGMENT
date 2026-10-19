package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotelmanagement/internal/domain"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("reservation conflicts with an existing reservation")
	ErrRoomUnavailable   = errors.New("room is out of service and cannot be booked")
	ErrInvalidTransition = errors.New("record cannot make that transition")
)

type tableSpec struct {
	table  string
	fields map[string]bool
}

var tables = map[domain.Kind]tableSpec{
	domain.Room:               {"rooms", fieldSet("number", "room_type", "floor", "status", "rate", "notes")},
	domain.Guest:              {"guests", fieldSet("first_name", "last_name", "email", "phone", "nationality", "document_number", "notes")},
	domain.Reservation:        {"reservations", fieldSet("room_id", "guest_id", "check_in", "check_out", "adults", "children", "rate", "source", "special_requests")},
	domain.HousekeepingTask:   {"housekeeping_tasks", fieldSet("room_id", "assigned_employee_id", "task_type", "priority", "status", "due_at", "notes")},
	domain.Folio:              {"folios", fieldSet("guest_id", "reservation_id", "status", "currency", "total_amount", "notes")},
	domain.Payment:            {"payments", fieldSet("folio_id", "amount", "currency", "method", "reference", "status")},
	domain.Employee:           {"employees", fieldSet("first_name", "last_name", "email", "phone", "role", "department", "status", "hired_on")},
	domain.MaintenanceRequest: {"maintenance_requests", fieldSet("room_id", "assigned_employee_id", "title", "description", "priority", "status", "reported_at", "resolved_at")},
	domain.InventoryItem:      {"inventory_items", fieldSet("name", "sku", "category", "unit", "quantity", "reorder_level", "unit_cost")},
	domain.InventoryMovement:  {"inventory_movements", fieldSet("item_id", "employee_id", "quantity", "movement_type", "note")},
}

func fieldSet(names ...string) map[string]bool {
	fields := make(map[string]bool, len(names))
	for _, name := range names {
		fields[name] = true
	}
	return fields
}

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context, kind domain.Kind, limit, offset int) ([]domain.Record, error) {
	spec := tables[kind]
	query := fmt.Sprintf(`SELECT id::text, '%s'::text, to_jsonb(t)-ARRAY['id','created_at','updated_at','archived_at'], created_at, updated_at
		FROM %s t WHERE archived_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, kind, spec.table)
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Record, 0)
	for rows.Next() {
		var item domain.Record
		if err := rows.Scan(&item.ID, &item.Kind, &item.Data, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *Repository) Get(ctx context.Context, kind domain.Kind, id string) (domain.Record, error) {
	spec := tables[kind]
	query := fmt.Sprintf(`SELECT id::text, '%s'::text, to_jsonb(t)-ARRAY['id','created_at','updated_at','archived_at'], created_at, updated_at
		FROM %s t WHERE id=$1 AND archived_at IS NULL`, kind, spec.table)
	return scanRecord(r.db.QueryRow(ctx, query, id))
}

func (r *Repository) Create(ctx context.Context, kind domain.Kind, data domain.Input) (domain.Record, error) {
	spec := tables[kind]
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return domain.Record{}, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return domain.Record{}, err
	}
	defer tx.Rollback(ctx)

	if kind == domain.Reservation {
		var roomID string
		if err := tx.QueryRow(ctx, `SELECT id::text FROM rooms WHERE id=$1 AND archived_at IS NULL FOR UPDATE`, data["room_id"]).Scan(&roomID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.Record{}, fmt.Errorf("room record not found")
			}
			return domain.Record{}, err
		}
		var outOfOrder bool
		if err := tx.QueryRow(ctx, `SELECT status='out_of_order' FROM rooms WHERE id=$1`, roomID).Scan(&outOfOrder); err != nil {
			return domain.Record{}, err
		}
		if outOfOrder {
			return domain.Record{}, ErrRoomUnavailable
		}
		var guestID string
		if err := tx.QueryRow(ctx, `SELECT id::text FROM guests WHERE id=$1 AND archived_at IS NULL`, data["guest_id"]).Scan(&guestID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.Record{}, fmt.Errorf("guest record not found")
			}
			return domain.Record{}, err
		}
		if data["status"] == nil || data["status"] == "" {
			data["status"] = "confirmed"
		}
	}
	columns, args := fieldsAndArgs(spec, data)
	placeholders := make([]string, len(args))
	for i := range args {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING id::text, '%s'::text,
		to_jsonb(%s)-ARRAY['id','created_at','updated_at','archived_at'], created_at, updated_at`, spec.table, strings.Join(columns, ","), strings.Join(placeholders, ","), kind, spec.table)
	item, err := scanRecord(tx.QueryRow(ctx, query, args...))
	if err != nil {
		return domain.Record{}, translateConstraint(err)
	}
	details, _ := json.Marshal(data)
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(record_kind,record_id,action,details) VALUES ($1,$2,'created',$3::jsonb)`, string(kind), item.ID, details); err != nil {
		return domain.Record{}, err
	}
	if kind == domain.InventoryMovement {
		if err := applyInventoryMovement(ctx, tx, data); err != nil {
			return domain.Record{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Record{}, translateConstraint(err)
	}
	return item, nil
}

func (r *Repository) Update(ctx context.Context, kind domain.Kind, id string, data domain.Input) (domain.Record, error) {
	spec := tables[kind]
	if len(data) == 0 {
		return r.Get(ctx, kind, id)
	}
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return domain.Record{}, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return domain.Record{}, err
	}
	defer tx.Rollback(ctx)
	if kind == domain.Reservation {
		var status string
		if err := tx.QueryRow(ctx, `SELECT status FROM reservations WHERE id=$1 AND archived_at IS NULL FOR UPDATE`, id).Scan(&status); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.Record{}, ErrNotFound
			}
			return domain.Record{}, err
		}
		if status != "confirmed" {
			return domain.Record{}, ErrInvalidTransition
		}
	}
	if kind == domain.Payment {
		var status string
		if err := tx.QueryRow(ctx, `SELECT status FROM payments WHERE id=$1 AND archived_at IS NULL FOR UPDATE`, id).Scan(&status); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.Record{}, ErrNotFound
			}
			return domain.Record{}, err
		}
		if status == "completed" && (data["amount"] != nil || data["folio_id"] != nil || data["method"] != nil) {
			return domain.Record{}, ErrInvalidTransition
		}
	}
	columns, args := fieldsAndArgs(spec, data)
	sets := make([]string, len(columns))
	for i, column := range columns {
		sets[i] = fmt.Sprintf("%s=$%d", column, i+1)
	}
	args = append(args, id)
	query := fmt.Sprintf(`UPDATE %s SET %s,updated_at=now() WHERE id=$%d AND archived_at IS NULL
		RETURNING id::text,'%s'::text,to_jsonb(%s)-ARRAY['id','created_at','updated_at','archived_at'],created_at,updated_at`, spec.table, strings.Join(sets, ","), len(args), kind, spec.table)
	item, err := scanRecord(tx.QueryRow(ctx, query, args...))
	if err != nil {
		return domain.Record{}, translateConstraint(err)
	}
	details, _ := json.Marshal(data)
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(record_kind,record_id,action,details) VALUES ($1,$2,'updated',$3::jsonb)`, string(kind), id, details); err != nil {
		return domain.Record{}, err
	}
	if kind == domain.HousekeepingTask {
		if status, ok := data["status"]; ok && status == "done" {
			if _, err := tx.Exec(ctx, `UPDATE rooms SET status='available',updated_at=now() WHERE id=(SELECT room_id FROM housekeeping_tasks WHERE id=$1) AND status='needs_cleaning'
				AND NOT EXISTS (SELECT 1 FROM housekeeping_tasks WHERE room_id=(SELECT room_id FROM housekeeping_tasks WHERE id=$1) AND id<>$1 AND archived_at IS NULL AND status IN ('open','in_progress'))`, id); err != nil {
				return domain.Record{}, err
			}
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Record{}, translateConstraint(err)
	}
	return item, nil
}

func (r *Repository) Archive(ctx context.Context, kind domain.Kind, id string) error {
	if kind == domain.Reservation || kind == domain.Payment || kind == domain.InventoryMovement {
		return fmt.Errorf("this record is part of an operational history and cannot be archived")
	}
	spec := tables[kind]
	query := fmt.Sprintf(`UPDATE %s SET archived_at=now(),updated_at=now() WHERE id=$1 AND archived_at IS NULL`, spec.table)
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	_, err = r.db.Exec(ctx, `INSERT INTO audit_events(record_kind,record_id,action) VALUES ($1,$2,'archived')`, string(kind), id)
	return err
}

func (r *Repository) TransitionReservation(ctx context.Context, id, action string) (domain.Record, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Record{}, err
	}
	defer tx.Rollback(ctx)
	var roomID, status string
	if err := tx.QueryRow(ctx, `SELECT room_id::text,status FROM reservations WHERE id=$1 AND archived_at IS NULL FOR UPDATE`, id).Scan(&roomID, &status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Record{}, ErrNotFound
		}
		return domain.Record{}, err
	}
	newStatus, roomStatus := "", ""
	switch action {
	case "check-in":
		if status != "confirmed" {
			return domain.Record{}, ErrInvalidTransition
		}
		var available bool
		if err := tx.QueryRow(ctx, `SELECT status='available' FROM rooms WHERE id=$1 AND archived_at IS NULL FOR UPDATE`, roomID).Scan(&available); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.Record{}, ErrNotFound
			}
			return domain.Record{}, err
		}
		if !available {
			return domain.Record{}, ErrInvalidTransition
		}
		newStatus, roomStatus = "checked_in", "occupied"
	case "check-out":
		if status != "checked_in" {
			return domain.Record{}, ErrInvalidTransition
		}
		newStatus, roomStatus = "checked_out", "needs_cleaning"
	case "cancel":
		if status != "confirmed" {
			return domain.Record{}, ErrInvalidTransition
		}
		newStatus = "cancelled"
	default:
		return domain.Record{}, ErrInvalidTransition
	}
	if _, err := tx.Exec(ctx, `UPDATE reservations SET status=$2,updated_at=now() WHERE id=$1`, id, newStatus); err != nil {
		return domain.Record{}, err
	}
	if roomStatus != "" {
		if _, err := tx.Exec(ctx, `UPDATE rooms SET status=$2,updated_at=now() WHERE id=$1`, roomID, roomStatus); err != nil {
			return domain.Record{}, err
		}
		if action == "check-out" {
			if _, err := tx.Exec(ctx, `INSERT INTO housekeeping_tasks(room_id,task_type,priority,status,notes) VALUES ($1,'checkout_cleaning','high','open','Auto-created at guest checkout')`, roomID); err != nil {
				return domain.Record{}, err
			}
		}
	}
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(record_kind,record_id,action) VALUES ('reservation',$1,$2)`, id, action); err != nil {
		return domain.Record{}, err
	}
	item, err := scanRecord(tx.QueryRow(ctx, `SELECT id::text,'reservation'::text,to_jsonb(t)-ARRAY['id','created_at','updated_at','archived_at'],created_at,updated_at FROM reservations t WHERE id=$1`, id))
	if err != nil {
		return domain.Record{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Record{}, err
	}
	return item, nil
}

func (r *Repository) Dashboard(ctx context.Context) (domain.Dashboard, error) {
	var d domain.Dashboard
	err := r.db.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM rooms WHERE archived_at IS NULL),
		(SELECT count(*) FROM rooms WHERE archived_at IS NULL AND status='available'),
		(SELECT count(*) FROM rooms WHERE archived_at IS NULL AND status='occupied'),
		(SELECT count(*) FROM reservations WHERE archived_at IS NULL AND check_in=current_date AND status='confirmed'),
		(SELECT count(*) FROM reservations WHERE archived_at IS NULL AND check_out=current_date AND status='checked_in'),
		(SELECT count(*) FROM housekeeping_tasks WHERE archived_at IS NULL AND status IN ('open','in_progress')),
		(SELECT count(*) FROM maintenance_requests WHERE archived_at IS NULL AND status IN ('open','in_progress')),
		(SELECT COALESCE(sum(amount),0) FROM payments WHERE archived_at IS NULL AND status='completed' AND created_at::date=current_date)`).Scan(
		&d.RoomsTotal, &d.RoomsAvailable, &d.RoomsOccupied, &d.ArrivalsToday,
		&d.DeparturesToday, &d.OpenTasks, &d.OpenMaintenance, &d.RevenueToday)
	return d, err
}

func (r *Repository) Occupancy(ctx context.Context) (domain.OccupancyReport, error) {
	var report domain.OccupancyReport
	err := r.db.QueryRow(ctx, `SELECT count(*),
		count(*) FILTER (WHERE status='available'),count(*) FILTER (WHERE status='occupied'),
		count(*) FILTER (WHERE status='needs_cleaning'),count(*) FILTER (WHERE status='out_of_order')
		FROM rooms WHERE archived_at IS NULL`).Scan(&report.Total, &report.Available, &report.Occupied, &report.NeedsCleaning, &report.OutOfOrder)
	if err == nil && report.Total > 0 {
		report.OccupancyPercent = float64(report.Occupied) * 100 / float64(report.Total)
	}
	return report, err
}

func (r *Repository) Revenue(ctx context.Context, from, to string) ([]domain.RevenueDay, error) {
	rows, err := r.db.Query(ctx, `SELECT created_at::date::text,count(*),sum(amount)::float8 FROM payments
		WHERE archived_at IS NULL AND status='completed' AND created_at::date >= $1::date AND created_at::date <= $2::date
		GROUP BY created_at::date ORDER BY created_at::date`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	report := make([]domain.RevenueDay, 0)
	for rows.Next() {
		var day domain.RevenueDay
		if err := rows.Scan(&day.Date, &day.Payments, &day.Revenue); err != nil {
			return nil, err
		}
		report = append(report, day)
	}
	return report, rows.Err()
}

func (r *Repository) FrontDeskQueue(ctx context.Context, direction, date string) ([]domain.FrontDeskEntry, error) {
	dateColumn, status := "r.check_in", "confirmed"
	if direction == "departures" {
		dateColumn, status = "r.check_out", "checked_in"
	}
	query := fmt.Sprintf(`SELECT r.id::text,r.guest_id::text,g.first_name || ' ' || g.last_name,
		r.room_id::text,rooms.number,r.check_in::text,r.check_out::text,r.status,r.adults,r.children
		FROM reservations r JOIN guests g ON g.id=r.guest_id JOIN rooms ON rooms.id=r.room_id
		WHERE r.archived_at IS NULL AND %s=COALESCE(NULLIF($1,'')::date,current_date) AND r.status=$2
		ORDER BY %s,rooms.number`, dateColumn, dateColumn)
	rows, err := r.db.Query(ctx, query, date, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := make([]domain.FrontDeskEntry, 0)
	for rows.Next() {
		var entry domain.FrontDeskEntry
		if err := rows.Scan(&entry.ReservationID, &entry.GuestID, &entry.GuestName, &entry.RoomID, &entry.RoomNumber, &entry.CheckIn, &entry.CheckOut, &entry.Status, &entry.Adults, &entry.Children); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func fieldsAndArgs(spec tableSpec, data domain.Input) ([]string, []any) {
	columns := make([]string, 0, len(data))
	for field := range data {
		if spec.fields[field] {
			columns = append(columns, field)
		}
	}
	sort.Strings(columns)
	args := make([]any, len(columns))
	for i, field := range columns {
		args[i] = data[field]
	}
	return columns, args
}

func applyInventoryMovement(ctx context.Context, tx pgx.Tx, data domain.Input) error {
	qty, ok := numeric(data["quantity"])
	if !ok || qty <= 0 {
		return fmt.Errorf("inventory movement quantity must be a positive number")
	}
	delta := qty
	switch data["movement_type"] {
	case "in", "restock":
	case "out", "use":
		delta = -qty
	default:
		return fmt.Errorf("movement_type must be in, out, restock, or use")
	}
	result, err := tx.Exec(ctx, `UPDATE inventory_items SET quantity=quantity+$2,updated_at=now()
		WHERE id=$1 AND archived_at IS NULL AND quantity+$2>=0`, data["item_id"], delta)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("inventory item not found or insufficient stock")
	}
	return nil
}

func numeric(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	case string:
		var f float64
		_, err := fmt.Sscan(v, &f)
		return f, err == nil
	default:
		return 0, false
	}
}

func scanRecord(row pgx.Row) (domain.Record, error) {
	var item domain.Record
	err := row.Scan(&item.ID, &item.Kind, &item.Data, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Record{}, ErrNotFound
	}
	return item, err
}

func translateConstraint(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23P01" {
			return ErrConflict
		}
		if pgErr.Code == "23505" {
			return fmt.Errorf("value already exists: %s", pgErr.ConstraintName)
		}
		if pgErr.Code == "23503" {
			return fmt.Errorf("referenced record does not exist")
		}
		if pgErr.Code == "23514" {
			return fmt.Errorf("value violates a hotel operations rule: %s", pgErr.ConstraintName)
		}
	}
	if strings.Contains(strings.ToLower(err.Error()), "invalid input syntax for type uuid") {
		return fmt.Errorf("record id must be a UUID")
	}
	return err
}
