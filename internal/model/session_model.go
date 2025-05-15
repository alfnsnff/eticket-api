package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ClaimSessionScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type ClaimSessionScheduleRoute struct {
	ID              uint                       `json:"id"`
	DepartureHarbor ClaimSessionScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   ClaimSessionScheduleHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type ClaimSessionScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// ScheduleDTO represents a Schedule.
type ClaimSessionSchedule struct {
	ID                uint                      `json:"id"`
	Ship              ClaimSessionScheduleShip  `json:"ship"`
	Route             ClaimSessionScheduleRoute `json:"route"`
	DepartureDatetime time.Time                 `json:"departure_datetime"`
	ArrivalDatetime   time.Time                 `json:"arrival_datetime"`
}

// ShipDTO represents a Ship.
type ReadClaimSessionResponse struct {
	ID          uint                                 `json:"id"`
	SessionID   string                               `json:"session_id"`
	ScheduleID  uint                                 `json:"schedule_id"`
	Schedule    ClaimSessionSchedule                 `json:"schedule"`
	ClaimedAt   time.Time                            `json:"claimed_at"`
	ExpiresAt   time.Time                            `json:"expires_at"`
	Prices      []ClaimSessionTicketPricesResponse   `json:"prices"`
	Tickets     []ClaimedSessionTicketDetailResponse `json:"tickets"`
	TotalAmount float32                              `json:"total_amount"`
	CreatedAt   time.Time                            `json:"created_at"`
	UpdatedAt   time.Time                            `json:"updated_at"`
}

type WriteClaimSessionRequest struct {
	SessionID   string    `json:"session_id"`
	ClaimedAt   time.Time `json:"claimed_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	TotalAmount float32   `json:"total_amount"`
}

type UpdateClaimSessionRequest struct {
	ID          uint      `json:"id"`
	SessionID   string    `json:"session_id"`
	ClaimedAt   time.Time `json:"claimed_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	TotalAmount float32   `json:"total_amount"`
}

type ClaimSessionTicketPricesResponse struct {
	ClassID  uint    `json:"class_id"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float32 `json:"subtotal"`
}

type ClaimedSessionTicketDetailResponse struct {
	TicketID uint                        `json:"ticket_id"`
	Class    ClaimSessionTicketClassItem `json:"class"`
	Price    float32                     `json:"price"` // Include price if frontend needs it at this stage
}

// ShipDTO represents a ship.
type ClaimSessionTicketClassItem struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
}

// ClaimTicketsRequest represents the input for claiming tickets
type ClaimedSessionLockTicketsRequest struct {
	ScheduleID uint                     `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimedSessionLockItem `json:"items"`       // List of classes and quantities requested
}

// ClaimItem represents a request for a specific class and quantity
type ClaimedSessionLockItem struct {
	ClassID  uint `json:"class_id"` // The class ID
	Quantity uint `json:"quantity"` // The number of tickets requested for this class
}

// ClaimTicketsResponse represents the result of a successful claim
type ClaimedSessionLockTicketsResponse struct {
	SessionID        string    `json:"session_id"`         // UUID for the session
	ClaimedTicketIDs []uint    `json:"claimed_ticket_ids"` // List of claimed ticket IDs
	ExpiresAt        time.Time `json:"expires_at"`         // Expiration time for the claim
}

type ClaimedSessionPassengerDataInput struct {
	TicketID      uint    `json:"ticket_id"`
	PassengerName string  `json:"passenger_name"`
	IDType        string  `json:"id_type"`
	IDNumber      string  `json:"id_number"`
	SeatNumber    *string `json:"seat_number"`
}

type ClaimedSessionFillPassengerDataRequest struct {
	SessionID     string                             `json:"session_id"`
	PassengerData []ClaimedSessionPassengerDataInput `json:"passenger_data"`
}

type ClaimedSessionFillPassengerDataResponse struct {
	UpdatedTicketIDs []uint                              `json:"updated_ticket_ids"`
	FailedTickets    []ClaimedSessionTicketUpdateFailure `json:"failed_tickets"`
}

type ClaimedSessionTicketUpdateFailure struct {
	TicketID uint   `json:"ticket_id"`
	Reason   string `json:"reason"`
	Test     []uint
}
