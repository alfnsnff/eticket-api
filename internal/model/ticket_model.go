package model

import (
	"time"
)

// TicketDTO represents a ticket.
type ReadTicketResponse struct {
	ID             uint      `json:"id"`
	ClaimSessionID uint      `json:"claim_session_id"`
	ScheduleID     uint      `json:"schedule_id"`
	ClassID        uint      `json:"class_id"`
	Status         string    `json:"status"`
	BookingID      uint      `json:"booking_id"`
	PassengerName  string    `json:"passenger_name"`
	PassengerAge   int       `json:"passenger_age"`
	Address        string    `json:"address"`
	SeatNumber     string    `json:"seat_number"`
	Price          float32   `json:"price"`
	IDType         string    `json:"id_type"`
	IDNumber       string    `json:"id_number"`
	ExpiresAt      time.Time `json:"expires_at"`
	ClaimedAt      time.Time `json:"claimed_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WriteTicketRequest struct {
	ClaimSessionID uint    `json:"claim_session_id"`
	ScheduleID     uint    `json:"schedule_id"`
	ClassID        uint    `json:"class_id"`
	Status         string  `json:"status"`
	BookingID      uint    `json:"booking_id"`
	PassengerName  string  `json:"passenger_name"`
	PassengerAge   int     `json:"passenger_age"`
	IDType         string  `json:"id_type"`
	IDNumber       string  `json:"id_number"`
	Address        string  `json:"address"`
	SeatNumber     string  `json:"seat_number"`
	Price          float32 `json:"price"`
}

type UpdateTicketRequest struct {
	ID             uint    `json:"id"`
	ClaimSessionID uint    `json:"claim_session_id"`
	ScheduleID     uint    `json:"schedule_id"`
	ClassID        uint    `json:"class_id"`
	Status         string  `json:"status"`
	BookingID      uint    `json:"booking_id"`
	PassengerName  string  `json:"passenger_name"`
	PassengerAge   int     `json:"passenger_age"`
	Address        string  `json:"address"`
	IDType         string  `json:"id_type"`
	IDNumber       string  `json:"id_number"`
	SeatNumber     string  `json:"seat_number"`
	Price          float32 `json:"price"`
}
