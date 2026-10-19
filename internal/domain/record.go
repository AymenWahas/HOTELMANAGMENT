package domain

import (
	"encoding/json"
	"time"
)

type Kind string

const (
	Room               Kind = "room"
	Guest              Kind = "guest"
	Reservation        Kind = "reservation"
	HousekeepingTask   Kind = "housekeeping_task"
	Folio              Kind = "folio"
	Payment            Kind = "payment"
	Employee           Kind = "employee"
	MaintenanceRequest Kind = "maintenance_request"
	InventoryItem      Kind = "inventory_item"
	InventoryMovement  Kind = "inventory_movement"
)

var Kinds = []Kind{Room, Guest, Reservation, HousekeepingTask, Folio, Payment, Employee, MaintenanceRequest, InventoryItem, InventoryMovement}

type Record struct {
	ID        string          `json:"id"`
	Kind      Kind            `json:"kind"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Input map[string]any

type Dashboard struct {
	RoomsTotal      int     `json:"rooms_total"`
	RoomsAvailable  int     `json:"rooms_available"`
	RoomsOccupied   int     `json:"rooms_occupied"`
	ArrivalsToday   int     `json:"arrivals_today"`
	DeparturesToday int     `json:"departures_today"`
	OpenTasks       int     `json:"open_housekeeping_tasks"`
	OpenMaintenance int     `json:"open_maintenance_requests"`
	RevenueToday    float64 `json:"revenue_today"`
}

type OccupancyReport struct {
	Total            int     `json:"total_rooms"`
	Available        int     `json:"available"`
	Occupied         int     `json:"occupied"`
	NeedsCleaning    int     `json:"needs_cleaning"`
	OutOfOrder       int     `json:"out_of_order"`
	OccupancyPercent float64 `json:"occupancy_percent"`
}

type RevenueDay struct {
	Date     string  `json:"date"`
	Payments int     `json:"payments"`
	Revenue  float64 `json:"revenue"`
}

type FrontDeskEntry struct {
	ReservationID string `json:"reservation_id"`
	GuestID       string `json:"guest_id"`
	GuestName     string `json:"guest_name"`
	RoomID        string `json:"room_id"`
	RoomNumber    string `json:"room_number"`
	CheckIn       string `json:"check_in"`
	CheckOut      string `json:"check_out"`
	Status        string `json:"status"`
	Adults        int    `json:"adults"`
	Children      int    `json:"children"`
}
