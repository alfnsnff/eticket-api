package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

type ShipClass struct {
	ID       uint `json:"id"`
	ShipID   uint `json:"ship_id"`
	ClassID  uint `json:"class_id"`
	Capacity int  `json:"capacity"`
}

// ShipDTO represents a Ship.
type ShipRead struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	ShipClasses []ShipClass `json:"ship_classes"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ShipDTO represents a Ship.
type ShipCreate struct {
	Name string `json:"name"`
}

func ToShipDTO(ship *entities.Ship) ShipRead {
	var shipResponse ShipRead
	copier.Copy(&shipResponse, &ship) // Automatically maps matching fields
	return shipResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToShipDTOs(ships []*entities.Ship) []ShipRead {
	var shipResponses []ShipRead
	for _, ship := range ships {
		shipResponses = append(shipResponses, ToShipDTO(ship))
	}
	return shipResponses
}

func ToShipEntity(shipCreate *ShipCreate) entities.Ship {
	var ship entities.Ship
	copier.Copy(&ship, &shipCreate)
	return ship
}
