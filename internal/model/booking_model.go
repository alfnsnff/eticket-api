package model

import (
	"time"
)

// HarborDTO represents a harbor.
type BookingScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type BookingScheduleRoute struct {
	ID              uint                  `json:"id"`
	DepartureHarbor BookingScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   BookingScheduleHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type BookingScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// ScheduleDTO represents a Schedule.
type BookingSchedule struct {
	ID                uint                 `json:"id"`
	Ship              BookingScheduleShip  `json:"ship"`
	Route             BookingScheduleRoute `json:"route"`
	DepartureDatetime time.Time            `json:"departure_datetime"`
	ArrivalDatetime   time.Time            `json:"arrival_datetime"`
}

// ShipDTO represents a ship.
type BookingTicketClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

// TicketDTO represents a ticket.
type BookingTicket struct {
	ID            uint               `json:"id"`
	Class         BookingTicketClass `json:"class"`
	Type          string             `json:"type" binding:"required,oneof=passenger vehicle"`
	PassengerName *string            `json:"passenger_name"`
	PassengerAge  *int               `json:"passenger_age"`
	Address       *string            `json:"address"`
	IDType        *string            `json:"id_type"`
	IDNumber      *string            `json:"id_number"`
	SeatNumber    *string            `json:"seat_number"`
	LicensePlate  *string            `json:"license_plate"`
	Price         float32            `json:"price"`
}

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID              uint            `json:"id"`
	OrderID         *string         `json:"order_id"` // Unique identifier for the booking, e.g., 'ORD123456'
	Schedule        BookingSchedule `json:"schedule"`
	CustomerName    string          `json:"customer_name"`
	CustomerAge     int             `json:"customer_age"`    // Age of the customer
	CUstomerGender  string          `json:"customer_gender"` //
	IDType          string          `json:"id_type"`
	IDNumber        string          `json:"id_number"`
	PhoneNumber     string          `json:"phone_number"` // Changed to string to support leading zeros
	Email           string          `json:"email"`
	ReferenceNumber *string         `json:"reference_number"` // Optional reference number for payment or external tracking
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Tickets         []BookingTicket `json:"tickets"`
}

type WriteBookingRequest struct {
	OrderID         *string `json:"order_id"`                                                    // Required
	ScheduleID      uint    `json:"schedule_id" validate:"required"`                             // Required
	IDType          string  `json:"id_type" validate:"required"`                                 // Required
	IDNumber        string  `json:"id_number" validate:"required"`                               // Required
	CustomerName    string  `json:"customer_name" validate:"required"`                           // Required
	CustomerAge     int     `json:"customer_age" validate:"required,min=0,max=120"`              // Basic age range
	CustomerGender  string  `json:"customer_gender" validate:"required,oneof=male female other"` // Must be one of these values
	PhoneNumber     string  `json:"phone_number" validate:"required"`                            // Uses E.164 format (e.g., +62812345678)
	Email           string  `json:"email" validate:"required,email"`                             // Must be valid email
	ReferenceNumber *string `json:"reference_number"`                                            // Optional
}

type UpdateBookingRequest struct {
	ID              uint    `json:"id,omitempty"`
	OrderID         *string `json:"order_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	CustomerName    string  `json:"customer_name" validate:"required"`
	CustomerAge     int     `json:"customer_age" validate:"required,min=0,max=120"`
	CustomerGender  string  `json:"customer_gender" validate:"required,oneof=male female other"`
	IDType          string  `json:"id_type" validate:"required"`
	IDNumber        string  `json:"id_number" validate:"required"`
	PhoneNumber     string  `json:"phone_number" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	ReferenceNumber *string `json:"reference_number"`
}
