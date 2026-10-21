package repository

import (
	"context"

	"hotelmanagement/internal/domain"
)

type RecordRepository interface {
	PublicRooms(context.Context, string, string) ([]domain.PublicRoom, error)
	CreateBookingRequest(context.Context, domain.BookingRequest) (domain.Record, error)
	List(context.Context, domain.Kind, int, int) ([]domain.Record, error)
	Get(context.Context, domain.Kind, string) (domain.Record, error)
	Create(context.Context, domain.Kind, domain.Input) (domain.Record, error)
	Update(context.Context, domain.Kind, string, domain.Input) (domain.Record, error)
	Archive(context.Context, domain.Kind, string) error
	TransitionReservation(context.Context, string, string) (domain.Record, error)
	Dashboard(context.Context) (domain.Dashboard, error)
	Occupancy(context.Context) (domain.OccupancyReport, error)
	Revenue(context.Context, string, string) ([]domain.RevenueDay, error)
	FrontDeskQueue(context.Context, string, string) ([]domain.FrontDeskEntry, error)
}
