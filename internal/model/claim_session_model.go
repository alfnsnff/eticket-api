package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ClaimSessionScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// ShipDTO represents a ship.
type ClaimSessionScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// ScheduleDTO represents a Schedule.
type ClaimSessionSchedule struct {
	ID                uint                       `json:"id"`
	Ship              ClaimSessionScheduleShip   `json:"ship"`
	DepartureHarbor   ClaimSessionScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor     ClaimSessionScheduleHarbor `json:"arrival_harbor"`
	DepartureDatetime time.Time                  `json:"departure_datetime"`
	ArrivalDatetime   time.Time                  `json:"arrival_datetime"`
}

// ShipDTO represents a ship.
type ClaimSessionTicketClassItem struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type ClaimSessionTicketPricesResponse struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`

	Class ClaimSessionTicketClassItem `json:"class"`
}

type ClaimSessionTicketDetailResponse struct {
	TicketID uint                        `json:"ticket_id"`
	Class    ClaimSessionTicketClassItem `json:"class"`
	Type     string                      `json:"type" binding:"required,oneof=passenger vehicle"`
	Price    float64                     `json:"price"` // Include price if frontend needs it at this stage
}

// ShipDTO represents a Ship.
type ReadClaimSessionResponse struct {
	ID          uint                               `json:"id"`
	SessionID   string                             `json:"session_id"`
	Schedule    ClaimSessionSchedule               `json:"schedule"`
	Status      string                             `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
	Tickets     []ClaimSessionTicketDetailResponse `json:"tickets"`
	Prices      []ClaimSessionTicketPricesResponse `json:"prices"`
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

// ClaimItem represents a request for a specific class and quantity
type ClaimSessionLockItem struct {
	ClassID  uint `json:"class_id"` // The class ID
	Quantity int  `json:"quantity"` // The number of tickets requested for this class
}

// ClaimTicketsRequest represents the input for claiming tickets
type WriteClaimSessionLockTicketsRequest struct {
	ScheduleID uint                   `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimSessionLockItem `json:"items"`       // List of classes and quantities requested
}

// ClaimTicketsResponse represents the result of a successful claim
type ReadClaimSessionLockTicketsResponse struct {
	SessionID        string    `json:"session_id"`       // UUID for the session
	ClaimedTicketIDs []uint    `json:"claim_ticket_ids"` // List of claim ticket IDs
	ExpiresAt        time.Time `json:"expires_at"`       // Expiration time for the claim
}

type ClaimSessionItem struct {
	ClassID  uint `json:"class_id"` // Foreign key to Class
	Quantity int  `json:"quantity"` // Number of tickets requested for this class
}

// ClaimTicketsResponse represents the result of a successful claim
type TESTReadClaimSessionLockTicketsResponse struct {
	ID          uint                               `json:"id"`
	SessionID   string                             `json:"session_id"`
	Schedule    ClaimSessionSchedule               `json:"schedule"`
	Status      string                             `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
	Tickets     []ClaimSessionTicketDetailResponse `json:"tickets"`
	Prices      []ClaimSessionTicketPricesResponse `json:"prices"`
	ClaimItems  []ClaimSessionItem                 `json:"claim_items"`
	TotalAmount float64                            `json:"total_amount"`
	ExpiresAt   time.Time                          `json:"expires_at"`
	CreatedAt   time.Time                          `json:"created_at"`
	UpdatedAt   time.Time                          `json:"updated_at"`
}

// ClaimTicketsRequest represents the input for claiming tickets
type TESTWriteClaimSessionLockTicketsRequest struct {
	ScheduleID uint               `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimSessionItem `json:"items"`       // List of classes and quantities requested
}

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
