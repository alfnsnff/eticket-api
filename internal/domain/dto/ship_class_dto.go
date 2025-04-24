package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// ShipDTO represents a Ship.
type ShipClassRead struct {
	ID        uint      `json:"id"`
	ShipID    uint      `json:"ship_id"`
	ClassID   uint      `json:"class_id"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ShipDTO represents a Ship.
type ShipClassCreate struct {
	ShipID   uint `json:"ship_id"`
	ClassID  uint `json:"class_id"`
	Capacity int  `json:"capacity"`
}

func ToShipClassDTO(shipClass *entities.ShipClass) ShipClassRead {
	var shipClassResponse ShipClassRead
	copier.Copy(&shipClassResponse, &shipClass) // Automatically maps matching fields
	return shipClassResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToShipClassDTOs(shipClasses []*entities.ShipClass) []ShipClassRead {
	var shipClassResponses []ShipClassRead
	for _, shipClass := range shipClasses {
		shipClassResponses = append(shipClassResponses, ToShipClassDTO(shipClass))
	}
	return shipClassResponses
}

func ToShipClassEntity(shipClassCreate *ShipClassCreate) entities.ShipClass {
	var shipClass entities.ShipClass
	copier.Copy(&shipClass, &shipClassCreate)
	return shipClass
}
