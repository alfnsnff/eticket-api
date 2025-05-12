package model

import (
	"time"
)

type CountFareRequest struct {
	ID         uint `json:"id"`
	RouteID    uint `json:"route_id"`
	ManifestID uint `json:"manifest_id"`
}

// ShipDTO represents a Ship.
type ReadFareResponse struct {
	ID          uint      `json:"id"`
	RouteID     uint      `json:"route_id"`
	ManifestID  uint      `json:"manifest_id"`
	TicketPrice float32   `json:"ticket_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WriteFareRequest struct {
	RouteID     uint    `json:"route_id"`
	ManifestID  uint    `json:"manifest_id"`
	TicketPrice float32 `json:"price"`
}

type UpdateFareRequest struct {
	ID          uint    `json:"id"`
	RouteID     uint    `json:"route_id"`
	ManifestID  uint    `json:"manifest_id"`
	TicketPrice float32 `json:"ticket_price"`
}
