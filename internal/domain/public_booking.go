package domain

type PublicRoom struct {
	ID       string  `json:"id"`
	RoomType string  `json:"room_type"`
	Rate     float64 `json:"rate"`
}

type BookingRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	RoomID    string `json:"room_id"`
	CheckIn   string `json:"check_in"`
	CheckOut  string `json:"check_out"`
	Adults    int    `json:"adults"`
	Children  int    `json:"children"`
	Website   string `json:"website"`
}
