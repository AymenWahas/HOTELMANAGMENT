package usecase

import (
	"context"
	"errors"
	"testing"

	"hotelmanagement/internal/domain"
	"hotelmanagement/internal/repository"
)

type bookingRequestRepository struct {
	repository.RecordRepository
	createdRequest *domain.BookingRequest
}

func (r *bookingRequestRepository) CreateBookingRequest(_ context.Context, request domain.BookingRequest) (domain.Record, error) {
	r.createdRequest = &request
	return domain.Record{ID: "booking-id"}, nil
}

func validBookingRequest() domain.BookingRequest {
	return domain.BookingRequest{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
		Phone:     "+1234567890",
		RoomID:    "room-id",
		CheckIn:   "2099-01-01",
		CheckOut:  "2099-01-02",
		Adults:    1,
	}
}

func TestRequestBookingAcceptsSinglePartName(t *testing.T) {
	repo := &bookingRequestRepository{}
	hotel := NewHotel(repo)
	request := validBookingRequest()
	request.FirstName = "  Madonna  "
	request.LastName = "   "

	if _, err := hotel.RequestBooking(context.Background(), request); err != nil {
		t.Fatalf("RequestBooking() error = %v, want nil", err)
	}
	if repo.createdRequest == nil {
		t.Fatal("CreateBookingRequest was not called")
	}
	if repo.createdRequest.FirstName != "Madonna" || repo.createdRequest.LastName != "" {
		t.Fatalf("stored name = %q %q, want %q and empty last name", repo.createdRequest.FirstName, repo.createdRequest.LastName, "Madonna")
	}
}

func TestRequestBookingRejectsMissingFirstName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
	}{
		{name: "missing"},
		{name: "whitespace only", firstName: " \t\n "},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &bookingRequestRepository{}
			hotel := NewHotel(repo)
			request := validBookingRequest()
			request.FirstName = test.firstName

			if _, err := hotel.RequestBooking(context.Background(), request); !errors.Is(err, ErrInvalid) {
				t.Fatalf("RequestBooking() error = %v, want ErrInvalid", err)
			}
			if repo.createdRequest != nil {
				t.Fatal("CreateBookingRequest was called for an invalid first name")
			}
		})
	}
}

func TestRequestBookingAcceptsTwoPartName(t *testing.T) {
	repo := &bookingRequestRepository{}
	hotel := NewHotel(repo)
	request := validBookingRequest()

	if _, err := hotel.RequestBooking(context.Background(), request); err != nil {
		t.Fatalf("RequestBooking() error = %v, want nil", err)
	}
	if repo.createdRequest == nil {
		t.Fatal("CreateBookingRequest was not called")
	}
	if repo.createdRequest.FirstName != "Ada" || repo.createdRequest.LastName != "Lovelace" {
		t.Fatalf("stored name = %q %q, want %q %q", repo.createdRequest.FirstName, repo.createdRequest.LastName, "Ada", "Lovelace")
	}
}

func TestRequestBookingPreservesNameLengthValidation(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
	}{
		{name: "first name too long", firstName: string(make([]byte, 81)), lastName: "Lovelace"},
		{name: "last name too long", firstName: "Ada", lastName: string(make([]byte, 81))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &bookingRequestRepository{}
			hotel := NewHotel(repo)
			request := validBookingRequest()
			request.FirstName = test.firstName
			request.LastName = test.lastName

			if _, err := hotel.RequestBooking(context.Background(), request); !errors.Is(err, ErrInvalid) {
				t.Fatalf("RequestBooking() error = %v, want ErrInvalid", err)
			}
			if repo.createdRequest != nil {
				t.Fatal("CreateBookingRequest was called for an overlong name")
			}
		})
	}
}
