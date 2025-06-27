package model

import (
	"time"
)

// HarborDTO represents a harbor.
type TicketScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// ShipDTO represents a ship.
type TicketScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// ScheduleDTO represents a Schedule.
type TicketSchedule struct {
	ID                uint                 `json:"id"`
	Ship              TicketScheduleShip   `json:"ship"`
	DepartureHarbor   TicketScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor     TicketScheduleHarbor `json:"arrival_harbor"`
	DepartureDatetime time.Time            `json:"departure_datetime"`
	ArrivalDatetime   time.Time            `json:"arrival_datetime"`
}

// ShipDTO represents a ship.
type TicketClassItem struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type TicketBooking struct {
	ID             uint   `json:"id"`
	OrderID        string `json:"order_id"` // Unique identifier for the booking, e.g., 'ORD123456'
	CustomerName   string `json:"customer_name"`
	CustomerAge    int    `json:"customer_age"`    // Age of the customer
	CUstomerGender string `json:"customer_gender"` //
	IDType         string `json:"id_type"`
	IDNumber       string `json:"id_number"`
	PhoneNumber    string `json:"phone_number"` // Changed to string to support leading zeros
	Email          string `json:"email"`
}

// TicketDTO represents a ticket.
type ReadTicketResponse struct {
	ID              uint            `json:"id"`
	Schedule        TicketSchedule  `json:"schedule"`
	Class           TicketClassItem `json:"class"`
	TicketCode      string          `json:"ticket_code"` // Unique ticket code
	BookingID       *uint           `json:"booking_id"`
	Booking         *TicketBooking  `json:"booking,omitempty"` // Optional booking details
	PassengerName   string          `json:"passenger_name"`
	PassengerAge    int             `json:"passenger_age"`
	Address         string          `json:"address"`
	PassengerGender *string         `json:"passenger"`
	IDType          *string         `json:"id_type"`
	IDNumber        *string         `json:"id_number"`
	SeatNumber      *string         `json:"seat_number"`
	LicensePlate    *string         `json:"license_plate"`
	Type            string          `json:"type" binding:"required,oneof=passenger vehicle"`
	Price           float64         `json:"price"`
	IsCheckedIn     bool            `json:"is_checked_in"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type WriteTicketRequest struct {
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	ClassID         uint    `json:"class_id" validate:"required"`
	BookingID       *uint   `json:"booking_id"`
	PassengerName   string  `json:"passenger_name"`
	PassengerAge    int     `json:"passenger_age"`
	Address         string  `json:"address"`
	PassengerGender *string `json:"passenger_gender"`
	IDType          *string `json:"id_type"`
	IDNumber        *string `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Type            string  `json:"type" validate:"required"`
	Price           float64 `json:"price" validate:"required,gte=0"`
	IsCheckedIn     bool    `json:"is_checked_in"`
}

type UpdateTicketRequest struct {
	ID              uint    `json:"id" validate:"required"`
	BookingID       *uint   `json:"booking_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	ClassID         uint    `json:"class_id" validate:"required"`
	PassengerName   string  `json:"passenger_name"`
	PassengerAge    int     `json:"passenger_age"`
	Address         string  `json:"address"`
	PassengerGender *string `json:"passenger_gender"`
	IDType          *string `json:"id_type"`
	IDNumber        *string `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Type            string  `json:"type" validate:"required"`
	Price           float64 `json:"price" validate:"required,gte=0"`
	IsCheckedIn     bool    `json:"is_checked_in"`
}
