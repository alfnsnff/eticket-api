package ship

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Ship entity to ReadShipResponse model
func ToReadShipResponse(ship *entity.Ship) *model.ReadShipResponse {
	if ship == nil {
		return nil
	}

	return &model.ReadShipResponse{
		ID:            ship.ID,
		ShipName:      ship.ShipName,
		Status:        ship.Status,
		ShipType:      ship.ShipType,
		ShipAlias:     ship.ShipAlias,
		YearOperation: ship.YearOperation,
		ImageLink:     ship.ImageLink,
		Description:   ship.Description,
		CreatedAt:     ship.CreatedAt,
		UpdatedAt:     ship.UpdatedAt,
	}
}

// Map a slice of Ship entities to ReadShipResponse models
func ToReadShipResponses(ships []*entity.Ship) []*model.ReadShipResponse {
	responses := make([]*model.ReadShipResponse, len(ships))
	for i, ship := range ships {
		responses[i] = ToReadShipResponse(ship)
	}
	return responses
}

// Map WriteShipRequest model to Ship entity
func FromWriteShipRequest(request *model.WriteShipRequest) *entity.Ship {
	return &entity.Ship{
		ShipName:      request.ShipName,
		Status:        request.Status,
		ShipType:      request.ShipType,
		ShipAlias:     request.ShipAlias,
		YearOperation: request.YearOperation,
		ImageLink:     request.ImageLink,
		Description:   request.Description,
	}
}

// Map UpdateShipRequest model to Ship entity
func FromUpdateShipRequest(request *model.UpdateShipRequest, ship *entity.Ship) {
	ship.ShipName = request.ShipName
	ship.Status = request.Status
	ship.ShipType = request.ShipType
	ship.ShipAlias = request.ShipAlias
	ship.YearOperation = request.YearOperation
	ship.ImageLink = request.ImageLink
	ship.Description = request.Description
}
