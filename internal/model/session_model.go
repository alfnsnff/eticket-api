package model

import (
	"time"
)

// ShipDTO represents a Ship.
type ReadClaimSessionResponse struct {
	ID          uint                               `json:"id"`
	SessionID   string                             `json:"session_id"`
	ClaimedAt   time.Time                          `json:"claimed_at"`
	ExpiresAt   time.Time                          `json:"expires_at"`
	Tickets     []ClaimSessionTicketPricesResponse `json:"tickets"`
	TotalAmount float32                            `json:"total_amount"`
	CreatedAt   time.Time                          `json:"created_at"`
	UpdatedAt   time.Time                          `json:"updated_at"`
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
