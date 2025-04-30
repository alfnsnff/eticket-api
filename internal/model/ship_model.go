package model

import (
	"time"
)

type ShipManifestClass struct {
	Name string `json:"name"`
}

type ShipManifest struct {
	ID       uint              `json:"id"`
	Class    ShipManifestClass `json:"class"`
	Capacity int               `json:"capacity"`
}

// ShipDTO represents a Ship.
type ReadShipResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Manifests []ShipManifest `json:"ship_classes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type WriteShipRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
