package model

import (
	"time"
)

type ManifestClass struct {
	Name string `json:"name"`
}

type ManifestShip struct {
	Name string `json:"name"`
}

// ShipDTO represents a Ship.
type ReadManifestResponse struct {
	ID        uint          `json:"id"`
	Class     ManifestClass `json:"class"`
	Ship      ManifestShip  `json:"ship"`
	Capacity  int           `json:"Manifest"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type WriteManifestRequest struct {
	ID       uint `json:"id"`
	ShipID   uint `json:"ship_id"`
	ClassID  uint `json:"class_id"`
	Manifest int  `json:"Manifest"`
}
