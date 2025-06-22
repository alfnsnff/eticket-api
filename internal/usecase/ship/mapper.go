package ship

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Ship domain to ReadShipResponse model
func ShipToResponse(ship *domain.Ship) *model.ReadShipResponse {
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
