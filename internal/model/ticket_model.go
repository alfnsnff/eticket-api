package model

import (
	"time"
)

// HarborDTO represents a harbor.
type TicketScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type TicketScheduleRoute struct {
	ID              uint                 `json:"id"`
	DepartureHarbor TicketScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   TicketScheduleHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type TicketScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// ScheduleDTO represents a Schedule.
type TicketSchedule struct {
	ID                uint                `json:"id"`
	Ship              TicketScheduleShip  `json:"ship"`
	Route             TicketScheduleRoute `json:"route"`
	DepartureDatetime time.Time           `json:"departure_datetime"`
	ArrivalDatetime   time.Time           `json:"arrival_datetime"`
}

// ShipDTO represents a ship.
type TicketClassItem struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

// TicketDTO represents a ticket.
type ReadTicketResponse struct {
	ID             uint            `json:"id"`
	ClaimSessionID uint            `json:"claim_session_id"`
	Schedule       TicketSchedule  `json:"schedule"`
	Class          TicketClassItem `json:"class"`
	Status         string          `json:"status"`
	BookingID      uint            `json:"booking_id"`
	Type           string          `json:"type" binding:"required,oneof=passenger vehicle"`
	PassengerName  string          `json:"passenger_name"`
	PassengerAge   int             `json:"passenger_age"`
	Address        string          `json:"address"`
	IDType         string          `json:"id_type"`
	IDNumber       string          `json:"id_number"`
	SeatNumber     *string         `json:"seat_number"`
	LicensePlate   *string         `json:"license_plate"`
	Price          float32         `json:"price"`
	ExpiresAt      time.Time       `json:"expires_at"`
	ClaimedAt      time.Time       `json:"claimed_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type WriteTicketRequest struct {
	ClaimSessionID  uint    `json:"claim_session_id"`
	ScheduleID      uint    `json:"schedule_id"`
	ClassID         uint    `json:"class_id"`
	Status          string  `json:"status"`
	BookingID       uint    `json:"booking_id"`
	Type            string  `json:"type" binding:"required,oneof=passenger vehicle"`
	PassengerName   string  `json:"passenger_name"`
	PassengerAge    int     `json:"passenger_age"`
	PassengerGender string  `json:"passenger_gender"`
	Address         string  `json:"address"`
	IDType          string  `json:"id_type"`
	IDNumber        string  `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Price           float32 `json:"price"`
}

type UpdateTicketRequest struct {
	ID              uint    `json:"id"`
	ClaimSessionID  uint    `json:"claim_session_id"`
	ScheduleID      uint    `json:"schedule_id"`
	ClassID         uint    `json:"class_id"`
	Status          string  `json:"status"`
	BookingID       uint    `json:"booking_id"`
	Type            string  `json:"type" binding:"required,oneof=passenger vehicle"`
	PassengerName   string  `json:"passenger_name"`
	PassengerAge    int     `json:"passenger_age"`
	PassengerGender string  `json:"passenger_gender"`
	Address         string  `json:"address"`
	IDType          string  `json:"id_type"`
	IDNumber        string  `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Price           float32 `json:"price"`
}
