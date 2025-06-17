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

// TicketDTO represents a ticket.
type BookingTicket struct {
	ID            uint            `json:"id"`
	Class         TicketClassItem `json:"class"`
	Status        string          `json:"status"`
	Type          string          `json:"type" binding:"required,oneof=passenger vehicle"`
	PassengerName string          `json:"passenger_name"`
	PassengerAge  int             `json:"passenger_age"`
	Address       string          `json:"address"`
	IDType        string          `json:"id_type"`
	IDNumber      string          `json:"id_number"`
	SeatNumber    *string         `json:"seat_number"`
	LicensePlate  *string         `json:"license_plate"`
	Price         float32         `json:"price"`
}

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID              uint            `json:"id"`
	OrderID         string          `json:"order_id"` // Unique identifier for the booking, e.g., 'ORD123456'
	Schedule        BookingSchedule `json:"schedule"`
	CustomerName    string          `json:"customer_name"`
	IDType          string          `json:"id_type"`
	IDNumber        string          `json:"id_number"`
	PhoneNumber     string          `json:"phone_number"` // Changed to string to support leading zeros
	Email           string          `json:"email_address"`
	Status          string          `json:"status"`           // e.g., 'completed', 'cancelled', 'refunded'
	ReferenceNumber *string         `json:"reference_number"` // Optional reference number for payment or external tracking

	BookedAt  time.Time `json:"booked_at"` // Timestamp when the booking was confirmed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tickets []BookingTicket `json:"tickets"`
}

type WriteBookingRequest struct {
	ScheduleID     uint   `json:"schedule_id"` // Foreign key
	IDType         string `json:"id_type"`
	IDNumber       string `json:"id_number"`
	CustomerName   string `json:"customer_name"`
	CustomerAge    int    `json:"customer_age"`
	CustomerGender string `json:"customer_gender"` //
	PhoneNumber    string `json:"phone_number"`    // Changed to string to support leading zeros
	Email          string `json:"email_address"`
	Status         string `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'
}

type UpdateBookingRequest struct {
	ID             uint   `json:"id"`
	ScheduleID     uint   `json:"schedule_id"` // Foreign key
	CustomerName   string `json:"customer_name"`
	CustomerAge    int    `json:"customer_age"`
	CustomerGender string `json:"customer_gender"` //
	IDType         string `json:"id_type"`
	IDNumber       string `json:"id_number"`
	PhoneNumber    string `json:"phone_number"` // Changed to string to support leading zeros
	Email          string `json:"email_address"`
	Status         string `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'
}

type ConfirmBookingRequest struct {
	Name        string    `json:"name"`
	IDType      string    `json:"id_type"`
	IDNumber    string    `json:"id_number"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	BirthDate   time.Time `json:"birth_date"`
	SessionID   string    `json:"session_id"`
	// TicketIDs   []uint    `json:"ticket_ids"` // List of ticket IDs being paid for
}

// ConfirmPaymentResponse represents the result of the payment confirmation.
type ConfirmBookingResponse struct {
	BookingID          uint   `json:"booking_id"`
	BookingStatus      string `json:"booking_status"`
	ConfirmedTicketIDs []uint `json:"confirmed_ticket_ids"`
}

type TicketSelectionResponse struct {
	ScheduleID uint                        `json:"schedule_id"`
	ShipName   string                      `json:"ship_name"`
	Datetime   string                      `json:"datetime"`
	Tickets    []TicketClassDetailResponse `json:"tickets"`
	Total      float32                     `json:"total"`
}

type TicketClassDetailResponse struct {
	ClassID  uint    `json:"class_id"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float32 `json:"subtotal"`
}
