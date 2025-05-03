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
