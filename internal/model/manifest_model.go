package model

import (
	"time"
)

// ShipDTO represents a Ship.
type ReadManifestResponse struct {
	ID        uint      `json:"id"`
	ClassID   uint      `json:"class_id"`
	ShipID    uint      `json:"ship_id"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteManifestRequest struct {
	ShipID   uint `json:"ship_id"`
	ClassID  uint `json:"class_id"`
	Capacity int  `json:"capacity"`
}

type UpdateManifestRequest struct {
	ID       uint `json:"id"`
	ShipID   uint `json:"ship_id"`
	ClassID  uint `json:"class_id"`
	Capacity int  `json:"capacity"`
}
