package model

import (
	"time"
)

// ShipDTO represents a Ship.
type ReadShipResponse struct {
	ID            uint      `json:"id"`
	ShipName      string    `json:"ship_name"`
	Status        string    `json:"status"`
	ShipType      string    `json:"ship_type"`
	ShipAlias     *string   `json:"ship_alias,omitempty"`
	YearOperation string    `json:"year_operation"`
	ImageLink     string    `json:"image_link"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WriteShipRequest struct {
	ShipName      string `json:"ship_name"`
	Status        string `json:"status"`
	ShipType      string `json:"ship_type"`
	YearOperation string `json:"year_operation"`
	ImageLink     string `json:"image_link"`
	Description   string `json:"description"`
}

type UpdateShipRequest struct {
	ID            uint   `json:"id"`
	ShipName      string `json:"ship_name"`
	Status        string `json:"status"`
	ShipType      string `json:"ship_type"`
	YearOperation string `json:"year_operation"`
	ImageLink     string `json:"image_link"`
	Description   string `json:"description"`
}
