package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"hotelmanagement/internal/domain"
	"hotelmanagement/internal/usecase"
)

type fakeRecordRepo struct {
	lastTransitionAction string
}

func (f *fakeRecordRepo) PublicRooms(context.Context, string, string) ([]domain.PublicRoom, error) {
	return nil, nil
}
func (f *fakeRecordRepo) CreateBookingRequest(context.Context, domain.BookingRequest) (domain.Record, error) {
	return domain.Record{ID: "booking-123"}, nil
}
func (f *fakeRecordRepo) List(context.Context, domain.Kind, int, int) ([]domain.Record, error) {
	return nil, nil
}
func (f *fakeRecordRepo) Get(context.Context, domain.Kind, string) (domain.Record, error) {
	return domain.Record{}, nil
}
func (f *fakeRecordRepo) Create(context.Context, domain.Kind, domain.Input) (domain.Record, error) {
	return domain.Record{}, nil
}
func (f *fakeRecordRepo) Update(context.Context, domain.Kind, string, domain.Input) (domain.Record, error) {
	return domain.Record{}, nil
}
func (f *fakeRecordRepo) Archive(context.Context, domain.Kind, string) error { return nil }
func (f *fakeRecordRepo) TransitionReservation(_ context.Context, _, action string) (domain.Record, error) {
	f.lastTransitionAction = action
	return domain.Record{ID: "booking-123", Kind: domain.Reservation, Data: json.RawMessage(`{"status":"cancelled"}`)}, nil
}
func (f *fakeRecordRepo) Dashboard(context.Context) (domain.Dashboard, error) {
	return domain.Dashboard{}, nil
}
func (f *fakeRecordRepo) Occupancy(context.Context) (domain.OccupancyReport, error) {
	return domain.OccupancyReport{}, nil
}
func (f *fakeRecordRepo) Revenue(context.Context, string, string) ([]domain.RevenueDay, error) {
	return nil, nil
}
func (f *fakeRecordRepo) FrontDeskQueue(context.Context, string, string) ([]domain.FrontDeskEntry, error) {
	return nil, nil
}

func TestPublicBookingRequestReturnsPendingStatus(t *testing.T) {
	bookingRequest := `{
		"first_name":"Ahmad",
		"last_name":"Ali",
		"email":"ahmad@example.com",
		"phone":"+966500000000",
		"room_id":"room-1",
		"check_in":"2099-01-10",
		"check_out":"2099-01-12",
		"adults":2,
		"children":0,
		"website":""
	}`

	h := NewHTTP(usecase.NewHotel(&fakeRecordRepo{}), "demo-token")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/booking-requests", strings.NewReader(bookingRequest))
	req.RemoteAddr = "127.0.0.1:5000"
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	h.publicBookingRequest(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status code = %d, want %d", res.Code, http.StatusCreated)
	}

	var payload map[string]string
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if payload["status"] != "pending" {
		t.Fatalf("status = %q, want %q", payload["status"], "pending")
	}
	if payload["id"] == "" {
		t.Fatal("expected booking id in response")
	}
	if _, err := time.Parse("2006-01-02", "2099-01-10"); err != nil {
		t.Fatalf("invalid test date: %v", err)
	}
}

func TestAllowBookingRequestAllowsReasonableVolume(t *testing.T) {
	h := NewHTTP(usecase.NewHotel(&fakeRecordRepo{}), "demo-token")
	addr := "198.51.100.15:443"

	for i := 0; i < 10; i++ {
		if !h.allowBookingRequest(addr) {
			t.Fatalf("request %d unexpectedly blocked for a normal guest", i+1)
		}
	}

	if h.allowBookingRequest(addr) {
		t.Fatal("11th attempt should be rejected after the configured per-window limit")
	}
}

func TestCancelBookingTransitionIsHandled(t *testing.T) {
	repo := &fakeRecordRepo{}
	h := NewHTTP(usecase.NewHotel(repo), "demo-token")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/booking-123/cancel", nil)
	req.SetPathValue("id", "booking-123")
	req.SetPathValue("action", "cancel")
	res := httptest.NewRecorder()

	h.reservationAction(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", res.Code, http.StatusOK)
	}
	if repo.lastTransitionAction != "cancel" {
		t.Fatalf("transition action = %q, want %q", repo.lastTransitionAction, "cancel")
	}
	if !strings.Contains(res.Body.String(), "cancelled") {
		t.Fatalf("cancel response body = %q, want cancelled state", res.Body.String())
	}
}
