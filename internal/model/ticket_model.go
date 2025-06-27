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

// TicketDTO represents a ticket.
type ReadTicketResponse struct {
	ID            uint            `json:"id"`
	Schedule      TicketSchedule  `json:"schedule"`
	Class         TicketClassItem `json:"class"`
	BookingID     *uint           `json:"booking_id"`
	PassengerName *string         `json:"passenger_name"`
	PassengerAge  *int            `json:"passenger_age"`
	Address       *string         `json:"address"`
	IDType        *string         `json:"id_type"`
	IDNumber      *string         `json:"id_number"`
	SeatNumber    *string         `json:"seat_number"`
	LicensePlate  *string         `json:"license_plate"`
	Type          string          `json:"type" binding:"required,oneof=passenger vehicle"`
	Price         float64         `json:"price"`
	IsCheckedIn   bool            `json:"is_checked_in"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type WriteTicketRequest struct {
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	ClassID         uint    `json:"class_id" validate:"required"`
	ClaimSessionID  *uint   `json:"claim_session_id"`
	BookingID       *uint   `json:"booking_id"`
	PassengerName   *string `json:"passenger_name"`
	PassengerAge    *int    `json:"passenger_age"`
	PassengerGender *string `json:"passenger_gender"`
	Address         *string `json:"address"`
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
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	ClassID         uint    `json:"class_id" validate:"required"`
	ClaimSessionID  *uint   `json:"claim_session_id"`
	BookingID       *uint   `json:"booking_id"`
	PassengerName   *string `json:"passenger_name"`
	PassengerAge    *int    `json:"passenger_age"`
	PassengerGender *string `json:"passenger_gender"`
	Address         *string `json:"address"`
	IDType          *string `json:"id_type"`
	IDNumber        *string `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Type            string  `json:"type" validate:"required"`
	Price           float64 `json:"price" validate:"required,gte=0"`
	IsCheckedIn     bool    `json:"is_checked_in"`
}
