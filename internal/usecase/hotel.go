package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"hotelmanagement/internal/domain"
	"hotelmanagement/internal/repository"
)

var ErrInvalid = errors.New("invalid input")

type Hotel struct{ records repository.RecordRepository }

func NewHotel(records repository.RecordRepository) *Hotel { return &Hotel{records: records} }

func (h *Hotel) List(ctx context.Context, kind domain.Kind, limit, offset int) ([]domain.Record, error) {
	if !validKind(kind) {
		return nil, fmt.Errorf("%w: unsupported resource", ErrInvalid)
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return h.records.List(ctx, kind, limit, offset)
}

func (h *Hotel) Get(ctx context.Context, kind domain.Kind, id string) (domain.Record, error) {
	if !validKind(kind) {
		return domain.Record{}, fmt.Errorf("%w: unsupported resource", ErrInvalid)
	}
	return h.records.Get(ctx, kind, id)
}

func (h *Hotel) Create(ctx context.Context, kind domain.Kind, data domain.Input) (domain.Record, error) {
	if !validKind(kind) {
		return domain.Record{}, fmt.Errorf("%w: unsupported resource", ErrInvalid)
	}
	if data == nil {
		return domain.Record{}, fmt.Errorf("%w: data is required", ErrInvalid)
	}
	if kind == domain.Guest {
		if email, ok := data["email"].(string); ok {
			email = strings.TrimSpace(email)
			if email == "" {
				delete(data, "email")
			} else {
				data["email"] = email
			}
		}
	}
	if err := validate(kind, data, true); err != nil {
		return domain.Record{}, err
	}
	return h.records.Create(ctx, kind, data)
}

func (h *Hotel) Update(ctx context.Context, kind domain.Kind, id string, data domain.Input) (domain.Record, error) {
	if !validKind(kind) {
		return domain.Record{}, fmt.Errorf("%w: unsupported resource", ErrInvalid)
	}
	if data == nil {
		return domain.Record{}, fmt.Errorf("%w: data is required", ErrInvalid)
	}
	if kind == domain.InventoryMovement {
		return domain.Record{}, fmt.Errorf("%w: inventory movements are immutable", ErrInvalid)
	}
	if err := validate(kind, data, false); err != nil {
		return domain.Record{}, err
	}
	return h.records.Update(ctx, kind, id, data)
}

func (h *Hotel) Archive(ctx context.Context, kind domain.Kind, id string) error {
	if !validKind(kind) {
		return fmt.Errorf("%w: unsupported resource", ErrInvalid)
	}
	if kind == domain.Reservation || kind == domain.Payment || kind == domain.InventoryMovement {
		return fmt.Errorf("%w: use the record's operational status instead of archiving history", ErrInvalid)
	}
	return h.records.Archive(ctx, kind, id)
}

func (h *Hotel) ReservationAction(ctx context.Context, id, action string) (domain.Record, error) {
	if action != "check-in" && action != "check-out" && action != "cancel" {
		return domain.Record{}, fmt.Errorf("%w: action must be check-in, check-out, or cancel", ErrInvalid)
	}
	return h.records.TransitionReservation(ctx, id, action)
}

func (h *Hotel) Dashboard(ctx context.Context) (domain.Dashboard, error) {
	return h.records.Dashboard(ctx)
}

func (h *Hotel) Occupancy(ctx context.Context) (domain.OccupancyReport, error) {
	return h.records.Occupancy(ctx)
}

func (h *Hotel) Revenue(ctx context.Context, from, to string) ([]domain.RevenueDay, error) {
	if from == "" || to == "" {
		return nil, fmt.Errorf("%w: from and to dates are required (YYYY-MM-DD)", ErrInvalid)
	}
	start, errStart := time.Parse("2006-01-02", from)
	end, errEnd := time.Parse("2006-01-02", to)
	if errStart != nil || errEnd != nil || end.Before(start) {
		return nil, fmt.Errorf("%w: from/to must be valid dates and to must not precede from", ErrInvalid)
	}
	if end.Sub(start) > 366*24*time.Hour {
		return nil, fmt.Errorf("%w: revenue report range cannot exceed one year", ErrInvalid)
	}
	return h.records.Revenue(ctx, from, to)
}

func (h *Hotel) FrontDeskQueue(ctx context.Context, direction, date string) ([]domain.FrontDeskEntry, error) {
	if direction != "arrivals" && direction != "departures" {
		return nil, fmt.Errorf("%w: reception queue must be arrivals or departures", ErrInvalid)
	}
	if date != "" {
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return nil, fmt.Errorf("%w: date must use YYYY-MM-DD", ErrInvalid)
		}
	}
	return h.records.FrontDeskQueue(ctx, direction, date)
}

func validKind(kind domain.Kind) bool {
	for _, allowed := range domain.Kinds {
		if kind == allowed {
			return true
		}
	}
	return false
}

var required = map[domain.Kind][]string{
	domain.Room:               []string{"number", "room_type"},
	domain.Guest:              []string{"first_name", "last_name"},
	domain.Reservation:        []string{"room_id", "guest_id", "check_in", "check_out"},
	domain.HousekeepingTask:   []string{"room_id", "task_type"},
	domain.Folio:              []string{"guest_id"},
	domain.Payment:            []string{"folio_id", "amount", "method"},
	domain.Employee:           []string{"first_name", "last_name", "role"},
	domain.MaintenanceRequest: []string{"title"},
	domain.InventoryItem:      []string{"name", "quantity"},
	domain.InventoryMovement:  []string{"item_id", "quantity", "movement_type"},
}

var allowedFields = map[domain.Kind]map[string]bool{
	domain.Room:               fields("number", "room_type", "floor", "status", "rate", "notes"),
	domain.Guest:              fields("first_name", "last_name", "email", "phone", "nationality", "document_number", "notes"),
	domain.Reservation:        fields("room_id", "guest_id", "check_in", "check_out", "adults", "children", "rate", "source", "special_requests"),
	domain.HousekeepingTask:   fields("room_id", "assigned_employee_id", "task_type", "priority", "status", "due_at", "notes"),
	domain.Folio:              fields("guest_id", "reservation_id", "status", "currency", "total_amount", "notes"),
	domain.Payment:            fields("folio_id", "amount", "currency", "method", "reference", "status"),
	domain.Employee:           fields("first_name", "last_name", "email", "phone", "role", "department", "status", "hired_on"),
	domain.MaintenanceRequest: fields("room_id", "assigned_employee_id", "title", "description", "priority", "status", "reported_at", "resolved_at"),
	domain.InventoryItem:      fields("name", "sku", "category", "unit", "quantity", "reorder_level", "unit_cost"),
	domain.InventoryMovement:  fields("item_id", "employee_id", "quantity", "movement_type", "note"),
}

func fields(names ...string) map[string]bool {
	result := make(map[string]bool, len(names))
	for _, name := range names {
		result[name] = true
	}
	return result
}

func validate(kind domain.Kind, data domain.Input, creating bool) error {
	for key := range data {
		if !allowedFields[kind][key] {
			return fmt.Errorf("%w: field %q is not supported for %s", ErrInvalid, key, kind)
		}
	}
	if creating {
		for _, key := range required[kind] {
			value, ok := data[key]
			if !ok || value == nil || strings.TrimSpace(fmt.Sprint(value)) == "" {
				return fmt.Errorf("%w: %s is required for %s", ErrInvalid, key, kind)
			}
		}
	}
	if kind == domain.Reservation {
		inValue, hasIn := data["check_in"]
		outValue, hasOut := data["check_out"]
		if hasIn && hasOut {
			inDate, errIn := time.Parse("2006-01-02", fmt.Sprint(inValue))
			outDate, errOut := time.Parse("2006-01-02", fmt.Sprint(outValue))
			if errIn != nil || errOut != nil || !outDate.After(inDate) {
				return fmt.Errorf("%w: check_in and check_out must be valid dates, and check_out must be later", ErrInvalid)
			}
		}
	}
	if kind == domain.Payment || kind == domain.InventoryItem || kind == domain.InventoryMovement {
		if amount, ok := data["amount"]; ok && (asFloat(amount) <= 0) {
			return fmt.Errorf("%w: amount must be greater than zero", ErrInvalid)
		}
		if quantity, ok := data["quantity"]; ok && (asFloat(quantity) < 0 || (kind == domain.InventoryMovement && asFloat(quantity) == 0)) {
			return fmt.Errorf("%w: quantity must be greater than zero for movements and cannot be negative for stock", ErrInvalid)
		}
	}
	return nil
}

func asFloat(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	default:
		var parsed float64
		_, _ = fmt.Sscan(fmt.Sprint(value), &parsed)
		return parsed
	}
}
