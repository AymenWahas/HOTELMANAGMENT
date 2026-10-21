package handler

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"hotelmanagement/internal/domain"
	"hotelmanagement/internal/postgres"
	"hotelmanagement/internal/usecase"
	"hotelmanagement/internal/webui"
)

type HTTP struct {
	hotel           *usecase.Hotel
	apiToken        string
	bookingMu       sync.Mutex
	bookingAttempts map[string][]time.Time
}

func NewHTTP(hotel *usecase.Hotel, apiToken string) *HTTP {
	return &HTTP{hotel: hotel, apiToken: apiToken, bookingAttempts: make(map[string][]time.Time)}
}

var resources = map[string]domain.Kind{
	"rooms": domain.Room, "guests": domain.Guest, "reservations": domain.Reservation, "bookings": domain.Reservation,
	"housekeeping": domain.HousekeepingTask, "folios": domain.Folio, "payments": domain.Payment,
	"employees": domain.Employee, "maintenance": domain.MaintenanceRequest,
	"inventory/items": domain.InventoryItem, "inventory/movements": domain.InventoryMovement,
}

func (h *HTTP) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/v1/public/rooms", h.publicRooms)
	mux.HandleFunc("POST /api/v1/public/booking-requests", h.publicBookingRequest)
	mux.HandleFunc("GET /api/v1/dashboard", h.dashboard)
	mux.HandleFunc("GET /api/v1/reports/occupancy", h.occupancy)
	mux.HandleFunc("GET /api/v1/reports/revenue", h.revenue)
	mux.HandleFunc("GET /api/v1/reception/arrivals", h.arrivals)
	mux.HandleFunc("GET /api/v1/reception/departures", h.departures)
	mux.HandleFunc("GET /api/v1/{resource...}", h.listOrGet)
	mux.HandleFunc("POST /api/v1/{resource...}", h.create)
	mux.HandleFunc("PATCH /api/v1/{resource...}", h.update)
	mux.HandleFunc("DELETE /api/v1/{resource...}", h.archive)
	mux.HandleFunc("POST /api/v1/reservations/{id}/{action}", h.reservationAction)
	mux.HandleFunc("POST /api/v1/bookings/{id}/{action}", h.reservationAction)
	mux.HandleFunc("POST /api/v1/reception/{id}/{action}", h.reservationAction)
	mux.Handle("/", webui.Handler())
	return cors(h.authenticate(mux))
}

func (h *HTTP) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPublicBookingRoute := (r.Method == http.MethodGet && r.URL.Path == "/api/v1/public/rooms") || (r.Method == http.MethodPost && r.URL.Path == "/api/v1/public/booking-requests")
		if strings.HasPrefix(r.URL.Path, "/api/") && !isPublicBookingRoute {
			provided := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if len(provided) != len(h.apiToken) || subtle.ConstantTimeCompare([]byte(provided), []byte(h.apiToken)) != 1 {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "valid bearer token required"})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h *HTTP) publicRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.hotel.PublicRooms(r.Context(), r.URL.Query().Get("check_in"), r.URL.Query().Get("check_out"))
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"items": rooms})
}

func (h *HTTP) publicBookingRequest(w http.ResponseWriter, r *http.Request) {
	if !h.allowBookingRequest(r.RemoteAddr) {
		w.Header().Set("Retry-After", "900")
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Too many booking requests. Please try again later."})
		return
	}
	var request domain.BookingRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, fmt.Errorf("%w: invalid JSON body", usecase.ErrInvalid))
		return
	}
	record, err := h.hotel.RequestBooking(r.Context(), request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": record.ID, "status": "pending"})
}

func (h *HTTP) allowBookingRequest(remoteAddr string) bool {
	clientIP, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		clientIP = remoteAddr
	}
	if clientIP == "" || clientIP == "[::]" || clientIP == "::" {
		return true
	}
	const maxAttemptsPerWindow = 10
	now := time.Now()
	cutoff := now.Add(-4 * time.Minute)
	h.bookingMu.Lock()
	defer h.bookingMu.Unlock()
	attempts := h.bookingAttempts[clientIP][:0]
	for _, attempt := range h.bookingAttempts[clientIP] {
		if attempt.After(cutoff) {
			attempts = append(attempts, attempt)
		}
	}
	if len(attempts) >= maxAttemptsPerWindow {
		h.bookingAttempts[clientIP] = attempts
		return false
	}
	h.bookingAttempts[clientIP] = append(attempts, now)
	return true
}

func (h *HTTP) dashboard(w http.ResponseWriter, r *http.Request) {
	result, err := h.hotel.Dashboard(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *HTTP) occupancy(w http.ResponseWriter, r *http.Request) {
	result, err := h.hotel.Occupancy(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *HTTP) revenue(w http.ResponseWriter, r *http.Request) {
	result, err := h.hotel.Revenue(r.Context(), r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *HTTP) arrivals(w http.ResponseWriter, r *http.Request) {
	h.frontDeskQueue(w, r, "arrivals")
}

func (h *HTTP) departures(w http.ResponseWriter, r *http.Request) {
	h.frontDeskQueue(w, r, "departures")
}

func (h *HTTP) frontDeskQueue(w http.ResponseWriter, r *http.Request, direction string) {
	entries, err := h.hotel.FrontDeskQueue(r.Context(), direction, r.URL.Query().Get("date"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": entries, "date": r.URL.Query().Get("date")})
}

func (h *HTTP) listOrGet(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	kind, ok := resources[resource]
	id := r.URL.Query().Get("id")
	if !ok {
		kind, id, ok = resourceAndID(resource)
		if !ok {
			writeError(w, usecase.ErrInvalid)
			return
		}
	}
	if id != "" {
		item, err := h.hotel.Get(r.Context(), kind, id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	items, err := h.hotel.List(r.Context(), kind, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "limit": limit, "offset": offset})
}

func (h *HTTP) create(w http.ResponseWriter, r *http.Request) {
	kind, ok := resources[r.PathValue("resource")]
	if !ok {
		writeError(w, usecase.ErrInvalid)
		return
	}
	var input domain.Input
	if err := decodeInput(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.hotel.Create(r.Context(), kind, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *HTTP) update(w http.ResponseWriter, r *http.Request) {
	kind, id, ok := resourceAndID(r.PathValue("resource"))
	if !ok {
		writeError(w, usecase.ErrInvalid)
		return
	}
	var input domain.Input
	if err := decodeInput(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.hotel.Update(r.Context(), kind, id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *HTTP) archive(w http.ResponseWriter, r *http.Request) {
	kind, id, ok := resourceAndID(r.PathValue("resource"))
	if !ok {
		writeError(w, usecase.ErrInvalid)
		return
	}
	if err := h.hotel.Archive(r.Context(), kind, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTP) reservationAction(w http.ResponseWriter, r *http.Request) {
	item, err := h.hotel.ReservationAction(r.Context(), r.PathValue("id"), r.PathValue("action"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func resourceAndID(path string) (domain.Kind, string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return "", "", false
	}
	resource := strings.Join(parts[:len(parts)-1], "/")
	kind, ok := resources[resource]
	if !ok {
		return "", "", false
	}
	return kind, parts[len(parts)-1], true
}

func decodeInput(r *http.Request, target *domain.Input) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"
	switch {
	case errors.Is(err, usecase.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, postgres.ErrNotFound):
		status, message = http.StatusNotFound, "record not found"
	case errors.Is(err, postgres.ErrConflict):
		status, message = http.StatusConflict, err.Error()
	case errors.Is(err, postgres.ErrRoomUnavailable):
		status, message = http.StatusConflict, err.Error()
	case errors.Is(err, postgres.ErrInvalidTransition):
		status, message = http.StatusConflict, err.Error()
	case strings.Contains(err.Error(), "record not found"):
		status, message = http.StatusNotFound, err.Error()
	case strings.Contains(err.Error(), "value already exists"), strings.Contains(err.Error(), "reservation conflicts"):
		status, message = http.StatusConflict, err.Error()
	case strings.Contains(err.Error(), "must be a UUID"), strings.Contains(err.Error(), "required"), strings.Contains(err.Error(), "must be"), strings.Contains(err.Error(), "invalid JSON body"), strings.Contains(err.Error(), "invalid date/time"), strings.Contains(err.Error(), "numeric value is out of range"), strings.Contains(err.Error(), "violates a hotel operations rule"), strings.Contains(err.Error(), "insufficient stock"), strings.Contains(err.Error(), "referenced record does not exist"):
		status, message = http.StatusBadRequest, err.Error()
	case strings.Contains(err.Error(), "cannot be archived"):
		status, message = http.StatusBadRequest, err.Error()
	default:
		message = "internal server error"
	}
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
