package model

import (
	"time"
)

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID           uint      `json:"id"`
	CustomerName string    `json:"customer_name"`
	IDType       uint      `json:"id_type"`
	IDNumber     uint      `json:"id_number"`
	PhoneNumber  string    `json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `json:"email_address"`
	BirthDate    time.Time `json:"birth_date"`

	BookingTimestamp time.Time `gorm:"not null" json:"booking_timestamp"`       // Timestamp when the booking was confirmed
	TotalAmount      float32   `gorm:"not null" json:"total_amount"`            // Total price of all tickets in this booking
	Status           string    `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteBookingRequest struct {
	ScheduleID   uint      `json:"schedule_id"` // Foreign key
	CustomerName string    `json:"customer_name"`
	PersonID     uint      `json:"person_id"`
	PhoneNumber  string    `json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `json:"email_address"`
	BirthDate    time.Time `json:"birth_date"`
}

type UpdateBookingRequest struct {
	ID           uint      `json:"id"`
	ScheduleID   uint      `json:"schedule_id"` // Foreign key
	CustomerName string    `json:"customer_name"`
	PersonID     uint      `json:"person_id"`
	PhoneNumber  string    `json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `json:"email_address"`
	BirthDate    time.Time `json:"birth_date"`
}

// ClaimTicketsRequest represents the input for claiming tickets
type LockTicketsRequest struct {
	ScheduleID uint       `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []LockItem `json:"items"`       // List of classes and quantities requested
}

// ClaimItem represents a request for a specific class and quantity
type LockItem struct {
	ClassID  uint `json:"class_id"` // The class ID
	Quantity uint `json:"quantity"` // The number of tickets requested for this class
}

// ClaimTicketsResponse represents the result of a successful claim
type LockTicketsResponse struct {
	SessionID        string    `json:"session_id"`         // UUID for the session
	ClaimedTicketIDs []uint    `json:"claimed_ticket_ids"` // List of claimed ticket IDs
	ExpiresAt        time.Time `json:"expires_at"`         // Expiration time for the claim
}

type ConfirmBookingRequest struct {
	Name        string    `json:"name"`
	IDType      string    `json:"id_type"`
	IDNumber    string    `json:"id_number"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	BirthDate   time.Time `json:"birth_date"`

	TicketIDs []uint `json:"ticket_ids"` // List of ticket IDs being paid for
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
