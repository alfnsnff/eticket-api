package model

import (
	"time"
)

type ManifestClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type ManifestShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
	ShipType string `json:"ship_type"`
}

// ShipDTO represents a Ship.
type ReadManifestResponse struct {
	ID        uint          `json:"id"`
	Class     ManifestClass `json:"class"`
	Ship      ManifestShip  `json:"ship"`
	Capacity  int           `json:"capacity"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
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
