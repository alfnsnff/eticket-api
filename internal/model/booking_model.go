package model

import (
	"eticket-api/pkg/utils/qr"
	"time"
)

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID           uint   `json:"id"`
	ScheduleID   uint   `json:"schedule_id"` // Foreign key
	CustomerName string `json:"customer_name"`
	IDType       uint   `json:"id_type"`
	IDNumber     uint   `json:"id_number"`
	PhoneNumber  string `json:"phone_number"` // Changed to string to support leading zeros
	Email        string `json:"email_address"`
	Status       string `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'

	BookedAt  time.Time `json:"booked_at"` // Timestamp when the booking was confirmed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteBookingRequest struct {
	ScheduleID   uint   `json:"schedule_id"` // Foreign key
	CustomerName string `json:"customer_name"`
	IDType       uint   `json:"id_type"`
	IDNumber     uint   `json:"id_number"`
	PhoneNumber  string `json:"phone_number"` // Changed to string to support leading zeros
	Email        string `json:"email_address"`
	Status       string `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'
}

type UpdateBookingRequest struct {
	ID           uint   `json:"id"`
	ScheduleID   uint   `json:"schedule_id"` // Foreign key
	CustomerName string `json:"customer_name"`
	IDType       uint   `json:"id_type"`
	IDNumber     uint   `json:"id_number"`
	PhoneNumber  string `json:"phone_number"` // Changed to string to support leading zeros
	Email        string `json:"email_address"`
	Status       string `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'
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
	BookingID          uint               `json:"booking_id"`
	BookingStatus      string             `json:"booking_status"`
	ConfirmedTicketIDs []uint             `json:"confirmed_ticket_ids"`
	Payment            qr.InvoiceResponse `json:"payment"` // QRIS payment details
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
