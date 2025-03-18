package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// ShipDTO represents a Ship.
type ShipRes struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShipName  string    `gorm:"not null" json:"ship_name"`
	Capacity  uint      `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToShipDTO(ship *entities.Ship) ShipRes {
	var shipResponse ShipRes
	copier.Copy(&shipResponse, &ship) // Automatically maps matching fields
	return shipResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToShipDTOs(ships []*entities.Ship) []ShipRes {
	var shipResponses []ShipRes
	for _, ship := range ships {
		shipResponses = append(shipResponses, ToShipDTO(ship))
	}
	return shipResponses
}
