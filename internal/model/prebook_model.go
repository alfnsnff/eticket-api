package model

import (
	"time"
)

// ClaimTicketsRequest represents the input for claiming tickets
type ClaimTicketsRequest struct {
	ScheduleID uint        `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimItem `json:"items"`       // List of classes and quantities requested
}

// ClaimItem represents a request for a specific class and quantity
type ClaimItem struct {
	ClassID  uint `json:"class_id"` // The class ID
	Quantity uint `json:"quantity"` // The number of tickets requested for this class
}

// ClaimTicketsResponse represents the result of a successful claim
type ClaimTicketsResponse struct {
	ClaimedTicketIDs []uint    `json:"claimed_ticket_ids"` // List of claimed ticket IDs
	ExpiresAt        time.Time `json:"expires_at"`         // Expiration time for the claim
}

// --- Input/Output Models ---

// PassengerDataInput represents data for a single passenger/ticket
type PassengerDataInput struct {
	TicketID      uint    `json:"ticket_id"`
	PassengerName string  `json:"passenger_name"`
	IDType        string  `json:"id_type"`     // e.g., Passport, ID Card
	IDNumber      string  `json:"id_number"`   // e.g., Passport Number, ID Card Number
	SeatNumber    *string `json:"seat_number"` // Pointer to string to allow NULL (if not assigned yet)
	// ... other passenger fields
}

// FillPassengerDataRequest represents the input for filling data for multiple tickets
type FillPassengerDataRequest struct {
	PassengerData []PassengerDataInput `json:"passenger_data"` // List of tickets and their data
}

// FillPassengerDataResponse represents the result of the data filling operation
type FillPassengerDataResponse struct {
	UpdatedTicketIDs []uint                `json:"updated_ticket_ids"`
	FailedTickets    []TicketUpdateFailure `json:"failed_tickets"` // Report tickets that failed to update
}

// TicketUpdateFailure details why a specific ticket update failed
type TicketUpdateFailure struct {
	TicketID uint   `json:"ticket_id"`
	Reason   string `json:"reason"`
}

type ConfirmPaymentRequest struct {
	// PaymentIntentID string `json:"payment_intent_id"` // Identifier from the payment gateway
	// Booker Information (Required fields since no user login)
	Name        string    `json:"name"`
	IDType      string    `json:"id_type"`
	IDNumber    string    `json:"id_number"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	BirthDate   time.Time `json:"birth_date"`

	TicketIDs []uint `json:"ticket_ids"` // List of ticket IDs being paid for
}

// ConfirmPaymentResponse represents the result of the payment confirmation.
type ConfirmPaymentResponse struct {
	BookingID          uint   `json:"booking_id"`
	BookingStatus      string `json:"booking_status"`
	ConfirmedTicketIDs []uint `json:"confirmed_ticket_ids"`
	// Potentially include details about the booking or confirmed tickets
}
