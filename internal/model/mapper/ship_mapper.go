package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToShipModel(ship *entities.Ship) *model.ReadShipResponse {
	response := new(model.ReadShipResponse)
	copier.Copy(&response, &ship) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToShipsModel(ships []*entities.Ship) []*model.ReadShipResponse {
	responses := []*model.ReadShipResponse{}
	for _, ship := range ships {
		responses = append(responses, ToShipModel(ship))
	}
	return responses
}

func ToShipEntity(request *model.WriteShipRequest) *entities.Ship {
	ship := new(entities.Ship)
	copier.Copy(&ship, &request)
	return ship
}
