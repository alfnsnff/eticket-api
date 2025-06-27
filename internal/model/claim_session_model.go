package model

import (
	"time"
)

// =========================
// Shared Components
// =========================
type ClaimSessionScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

type ClaimSessionScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

type ClaimSessionSchedule struct {
	ID                uint                       `json:"id"`
	Ship              ClaimSessionScheduleShip   `json:"ship"`
	DepartureHarbor   ClaimSessionScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor     ClaimSessionScheduleHarbor `json:"arrival_harbor"`
	DepartureDatetime time.Time                  `json:"departure_datetime"`
	ArrivalDatetime   time.Time                  `json:"arrival_datetime"`
}

type ClaimSessionItemClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type ClaimSessionItem struct {
	ClassID  uint                  `json:"class_id"`
	Class    ClaimSessionItemClass `json:"class"`
	Quantity int                   `json:"quantity"`
}

type ClaimSessionTicket struct {
	TicketID uint                  `json:"ticket_id"`
	Class    ClaimSessionItemClass `json:"class"`
	Type     string                `json:"type" binding:"required,oneof=passenger vehicle"`
	Price    float64               `json:"price"`
}

type ClaimSessionTicketPricesResponse struct {
	Price    float64               `json:"price"`
	Quantity int                   `json:"quantity"`
	Class    ClaimSessionItemClass `json:"class"`
	Subtotal float64               `json:"subtotal"`
}

// =========================
// ClaimSession Main Flow
// =========================
type ReadClaimSessionResponse struct {
	ID          uint                               `json:"id"`
	SessionID   string                             `json:"session_id"`
	Schedule    ClaimSessionSchedule               `json:"schedule"`
	Status      string                             `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
	Tickets     []ClaimSessionTicket               `json:"tickets"`
	Prices      []ClaimSessionTicketPricesResponse `json:"prices"`
	ClaimItems  []ClaimSessionItem                 `json:"claim_items"`
	TotalAmount float64                            `json:"total_amount"`
	ExpiresAt   time.Time                          `json:"expires_at"`
	CreatedAt   time.Time                          `json:"created_at"`
	UpdatedAt   time.Time                          `json:"updated_at"`
}

type WriteClaimSessionRequest struct {
	SessionID   string    `json:"session_id"`
	ClaimAt     time.Time `json:"claim_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	TotalAmount float32   `json:"total_amount"`
}

type UpdateClaimSessionRequest struct {
	ID          uint      `json:"id"`
	SessionID   string    `json:"session_id"`
	ClaimAt     time.Time `json:"claim_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	TotalAmount float32   `json:"total_amount"`
}

// =========================
// ClaimSession Lock Tickets Flow
// =========================
type ClaimSessionLockItem struct {
	ClassID  uint `json:"class_id"`
	Quantity int  `json:"quantity"` // The number of tickets requested for this class
}

type WriteClaimSessionLockTicketsRequest struct {
	ScheduleID uint                   `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimSessionLockItem `json:"items"`       // List of classes and quantities requested
}

type ReadClaimSessionLockTicketsResponse struct {
	SessionID        string    `json:"session_id"`       // UUID for the session
	ClaimedTicketIDs []uint    `json:"claim_ticket_ids"` // List of claim ticket IDs
	ExpiresAt        time.Time `json:"expires_at"`       // Expiration time for the claim
}

type TESTReadClaimSessionLockResponse struct {
	SessionID string    `json:"session_id"` // UUID for the session
	ExpiresAt time.Time `json:"expires_at"` // Expiration time for the claim
}

// Extended/test version with full session details
type TESTReadClaimSessionResponse struct {
	ID        uint                 `json:"id"`
	SessionID string               `json:"session_id"`
	Schedule  ClaimSessionSchedule `json:"schedule"`
	Status    string               `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
	// Tickets     []ClaimSessionTicket               `json:"tickets"`
	Prices      []ClaimSessionTicketPricesResponse `json:"prices"`
	Tickets     []ClaimSessionItem                 `json:"claim_items"`
	TotalAmount float64                            `json:"total_amount"`
	ExpiresAt   time.Time                          `json:"expires_at"`
	CreatedAt   time.Time                          `json:"created_at"`
	UpdatedAt   time.Time                          `json:"updated_at"`
}

type TESTWriteClaimSessionRequest struct {
	ScheduleID uint               `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimSessionItem `json:"items"`       // List of classes and quantities requested
}

type TESTClaimSessionTicketDataEntry struct {
	ClassID         uint    `json:"class_id"` // The class ID for the ticket
	PassengerName   string  `json:"passenger_name"`
	IDType          string  `json:"id_type"`
	IDNumber        string  `json:"id_number"`
	PassengerAge    int     `json:"passenger_age"`
	PassengerGender string  `json:"passenger_gender"`
	Address         string  `json:"address"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
}

type TESTWriteClaimSessionDataEntryRequest struct {
	CustomerName   string                            `json:"customer_name"`
	CustomerAge    int                               `json:"customer_age"`
	CustomerGender string                            `json:"customer_gender"`
	IDType         string                            `json:"id_type"`
	IDNumber       string                            `json:"id_number"`
	PhoneNumber    string                            `json:"phone_number"`
	Email          string                            `json:"email"`
	BirthDate      time.Time                         `json:"birth_date"`
	TicketData     []TESTClaimSessionTicketDataEntry `json:"ticket_data"`
	PaymentMethod  string                            `json:"payment_method"`
}

// =========================
// ClaimSession Data Entry Flow
// =========================
type ClaimSessionTicketDataEntry struct {
	TicketID        uint    `json:"ticket_id"`
	PassengerName   string  `json:"passenger_name"`
	IDType          string  `json:"id_type"`
	IDNumber        string  `json:"id_number"`
	PassengerAge    int     `json:"passenger_age"`
	PassengerGender string  `json:"passenger_gender"`
	Address         string  `json:"address"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
}

type WriteClaimSessionDataEntryRequest struct {
	SessionID      string                        `json:"session_id"`
	CustomerName   string                        `json:"customer_name"`
	CustomerAge    int                           `json:"customer_age"`
	CustomerGender string                        `json:"customer_gender"`
	IDType         string                        `json:"id_type"`
	IDNumber       string                        `json:"id_number"`
	PhoneNumber    string                        `json:"phone_number"`
	Email          string                        `json:"email"`
	BirthDate      time.Time                     `json:"birth_date"`
	TicketData     []ClaimSessionTicketDataEntry `json:"ticket_data"`
}

type ReadClaimSessionDataEntryResponse struct {
	BookingID        uint   `json:"booking_id"` // ID of the booking created
	OrderID          string `json:"order_id"`   // ID of the booking created
	UpdatedTicketIDs []uint `json:"updated_ticket_ids"`
}
