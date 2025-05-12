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
	SeatNumber     string    `json:"seat_number"`
	Price          float32   `json:"price"`
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
	SeatNumber     string  `json:"seat_number"`
	Price          float32 `json:"price"`
}

type PassengerDataInput struct {
	TicketID      uint    `json:"ticket_id"`
	PassengerName string  `json:"passenger_name"`
	IDType        string  `json:"id_type"`
	IDNumber      string  `json:"id_number"`
	SeatNumber    *string `json:"seat_number"`
}

type FillPassengerDataRequest struct {
	SessionID     string               `json:"session_id"`
	PassengerData []PassengerDataInput `json:"passenger_data"`
}

type FillPassengerDataResponse struct {
	UpdatedTicketIDs []uint                `json:"updated_ticket_ids"`
	FailedTickets    []TicketUpdateFailure `json:"failed_tickets"`
}

type TicketUpdateFailure struct {
	TicketID uint   `json:"ticket_id"`
	Reason   string `json:"reason"`
	Test     []uint
}
