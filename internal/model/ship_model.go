package model

import (
	"time"
)

// type ShipManifestClass struct {
// 	Name string `json:"name"`
// }

// type ShipManifest struct {
// 	ID       uint              `json:"id"`
// 	Class    ShipManifestClass `json:"class"`
// 	Capacity int               `json:"capacity"`
// }

// ShipDTO represents a Ship.
type ReadShipResponse struct {
	ID          uint      `json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Status      string    `json:"status"`
	ShipType    string    `json:"ship_type"`
	Year        string    `json:"year"`
	Image       string    `json:"image"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WriteShipRequest struct {
	Name        string `gorm:"not null" json:"name"`
	Status      string `json:"status"`
	ShipType    string `json:"ship_type"`
	Year        string `json:"year"`
	Image       string `json:"image"`
	Description string `json:"Description"`
}

type UpdateShipRequest struct {
	ID          uint   `json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Status      string `json:"status"`
	ShipType    string `json:"ship_type"`
	Year        string `json:"year"`
	Image       string `json:"image"`
	Description string `json:"Description"`
}
