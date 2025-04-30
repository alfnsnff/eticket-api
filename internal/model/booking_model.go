package model

import (
	"time"
)

// HarborDTO represents a harbor.
type BookingHarbor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RouteDTO represents a travel route.
type BookingRoute struct {
	ID              uint          `json:"id"`
	DepartureHarbor BookingHarbor `json:"departure_harbor"`
	ArrivalHarbor   BookingHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type BookingShip struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ScheduleDTO represents a trip schedule.
type BookingSchedule struct {
	ID       uint         `json:"id"`
	DateTime time.Time    `json:"datetime"`
	Ship     BookingShip  `json:"ship"`
	Route    BookingRoute `json:"route"`
}

// ClassDTO represents ticket class information.
type BookingTicketClass struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookingShipClass struct {
	ID    uint        `json:"id"`
	Class TicketClass `json:"class"`
}

type BookingTicketFare struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Price     float32          `json:"price"`
	ShipClass BookingShipClass `json:"ship_class"`
}

type CreateBookingTicketRequest struct {
	PriceID       uint   `json:"price_id"`
	ScheduleID    uint   `json:"schedule_id"`
	PassengerName string `json:"passenger_name"`
	SeatNumber    string `json:"seat_number"`
}

type ReadBookingTicketResponse struct {
	PassengerName string            `json:"passenger_name"`
	SeatNumber    string            `json:"seat_number"`
	Fare          BookingTicketFare `json:"fare"`
}

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID          uint                        `json:"id"`
	CusName     string                      `json:"cus_name"`
	PersonID    uint                        `json:"person_id"`
	PhoneNumber string                      `json:"phone_number"` // Changed to string to support leading zeros
	Email       string                      `json:"email_address"`
	BirthDate   time.Time                   `json:"birth_date"`
	Schedule    BookingSchedule             `json:"schedule"`
	Tickets     []ReadBookingTicketResponse `json:"tickets"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type WriteBookingRequest struct {
	ID          uint                         `json:"id"`
	ScheduleID  uint                         `json:"schedule_id"` // Foreign key
	CusName     string                       `json:"cus_name"`
	PersonID    uint                         `json:"person_id"`
	PhoneNumber string                       `json:"phone_number"` // Changed to string to support leading zeros
	Email       string                       `json:"email_address"`
	BirthDate   time.Time                    `json:"birth_date"`
	Schedule    BookingSchedule              `json:"schedule"`
	Tickets     []CreateBookingTicketRequest `json:"tickets"`
}
