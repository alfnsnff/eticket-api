package model

import (
	"time"
)

// // HarborDTO represents a harbor.
// type TicketHarbor struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // RouteDTO represents a travel route.
// type TicketRoute struct {
// 	ID              uint         `json:"id"`
// 	DepartureHarbor TicketHarbor `json:"departure_harbor"`
// 	ArrivalHarbor   TicketHarbor `json:"arrival_harbor"`
// }

// // ShipDTO represents a ship.
// type TicketShip struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // ScheduleDTO represents a trip schedule.
// type TicketSchedule struct {
// 	ID       uint        `json:"id"`
// 	DateTime time.Time   `json:"datetime"`
// 	Ship     TicketShip  `json:"ship"`
// 	Route    TicketRoute `json:"route"`
// }

// // ClassDTO represents ticket class information.
// type TicketClass struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// type TicketManifest struct {
// 	ID    uint        `json:"id"`
// 	Class TicketClass `json:"class"`
// }

// type TicketFare struct {
// 	ID       uint           `json:"id"`
// 	Price    float32        `json:"price"`
// 	Manifest TicketManifest `json:"manifest"`
// }

// // BookingDTO represents the person who booked the ticket.
// type TicketBooking struct {
// 	ID          uint      `json:"id"`
// 	CusName     string    `json:"cus_name"`
// 	PersonID    uint      `json:"person_id"`
// 	PhoneNumber string    `json:"phone_number"` // Changed to string to support leading zeros
// 	Email       string    `json:"email_address"`
// 	BirthDate   time.Time `json:"birth_date"`
// }

// TicketDTO represents a ticket.
type ReadTicketResponse struct {
	ID            uint      `json:"id"`
	ScheduleID    uint      `json:"schedule_id"`
	ClassID       uint      `json:"class_id"`
	Status        string    `json:"status"`
	BookingID     uint      `json:"booking_id"`
	PassengerName string    `json:"passenger_name"`
	SeatNumber    string    `json:"seat_number"`
	Price         float32   `json:"price"`
	ExpiresAt     time.Time `json:"expires_at"` // Added not null, essential for timeout
	ClaimedAt     time.Time `json:"claimed_at"` // Added not null, essential for timeout
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WriteTicketRequest struct {
	ScheduleID    uint    `json:"schedule_id"`
	ClassID       uint    `json:"class_id"`
	Status        string  `json:"status"`
	BookingID     uint    `json:"booking_id"`
	PassengerName string  `json:"passenger_name"`
	SeatNumber    string  `json:"seat_number"`
	Price         float32 `json:"price"`
}

type UpdateTicketRequest struct {
	ID            uint    `json:"id"`
	ScheduleID    uint    `json:"schedule_id"`
	ClassID       uint    `json:"class_id"`
	Status        string  `json:"status"`
	BookingID     uint    `json:"booking_id"`
	PassengerName string  `json:"passenger_name"`
	SeatNumber    string  `json:"seat_number"`
	Price         float32 `json:"price"`
}

// type CountBookedTicketRequest struct {
// 	ID         uint `json:"id"`
// 	ScheduleID uint `json:"schedule_id"`
// 	FareID     uint `json:"fare_id"`
// }

// type TicketSelectionRequest struct {
// 	ScheduleID uint                  `json:"schedule_id" binding:"required"`
// 	Tickets    []TicketClassQuantity `json:"tickets" binding:"required"`
// }

// type TicketClassQuantity struct {
// 	FareID   uint `json:"fare_id" binding:"required"`
// 	Quantity int  `json:"quantity" binding:"required,min=1"`
// }

// type TicketSelectionResponse struct {
// 	ScheduleID uint                        `json:"schedule_id"`
// 	ShipName   string                      `json:"ship_name"`
// 	Datetime   string                      `json:"datetime"`
// 	Tickets    []TicketClassDetailResponse `json:"tickets"`
// 	Total      float32                     `json:"total"`
// }

// type TicketClassDetailResponse struct {
// 	ClassName string  `json:"class_name"`
// 	FareID    uint    `json:"fare_id"`
// 	Price     float32 `json:"price"`
// 	Quantity  int     `json:"quantity"`
// 	Subtotal  float32 `json:"subtotal"`
// }
