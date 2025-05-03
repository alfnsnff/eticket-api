package model

import (
	"time"
)

// // HarborDTO represents a harbor.
// type BookingHarbor struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // RouteDTO represents a travel route.
// type BookingRoute struct {
// 	ID              uint          `json:"id"`
// 	DepartureHarbor BookingHarbor `json:"departure_harbor"`
// 	ArrivalHarbor   BookingHarbor `json:"arrival_harbor"`
// }

// // ShipDTO represents a ship.
// type BookingShip struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // ScheduleDTO represents a trip schedule.
// type BookingSchedule struct {
// 	ID       uint         `json:"id"`
// 	DateTime time.Time    `json:"datetime"`
// 	Ship     BookingShip  `json:"ship"`
// 	Route    BookingRoute `json:"route"`
// }

// // ClassDTO represents ticket class information.
// type BookingTicketClass struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// type BookingManifest struct {
// 	ID    uint        `json:"id"`
// 	Class TicketClass `json:"class"`
// }

// type BookingTicketFare struct {
// 	ID       uint            `json:"id"`
// 	Price    float32         `json:"price"`
// 	Manifest BookingManifest `json:"manifest"`
// }

// type WriteBookingTicketRequest struct {
// 	FareID        uint   `json:"fare_id"`
// 	ScheduleID    uint   `json:"schedule_id"`
// 	PassengerName string `json:"passenger_name"`
// 	SeatNumber    string `json:"seat_number"`
// }

// type ReadBookingTicketResponse struct {
// 	PassengerName string            `json:"passenger_name"`
// 	SeatNumber    string            `json:"seat_number"`
// 	Fare          BookingTicketFare `json:"fare"`
// }

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID          uint      `json:"id"`
	CusName     string    `json:"cus_name"`
	PersonID    uint      `json:"person_id"`
	PhoneNumber string    `json:"phone_number"` // Changed to string to support leading zeros
	Email       string    `json:"email_address"`
	BirthDate   time.Time `json:"birth_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WriteBookingRequest struct {
	SessionID   string    `json:"session_id"`  // UUID from client or backend
	ScheduleID  uint      `json:"schedule_id"` // Foreign key
	CusName     string    `json:"cus_name"`
	PersonID    uint      `json:"person_id"`
	PhoneNumber string    `json:"phone_number"` // Changed to string to support leading zeros
	Email       string    `json:"email_address"`
	BirthDate   time.Time `json:"birth_date"`
}

type UpdateBookingRequest struct {
	ID          uint      `json:"id"`
	SessionID   string    `json:"session_id"`  // UUID from client or backend
	ScheduleID  uint      `json:"schedule_id"` // Foreign key
	CusName     string    `json:"cus_name"`
	PersonID    uint      `json:"person_id"`
	PhoneNumber string    `json:"phone_number"` // Changed to string to support leading zeros
	Email       string    `json:"email_address"`
	BirthDate   time.Time `json:"birth_date"`
}
