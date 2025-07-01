package requests

import (
	"eticket-api/internal/domain"
	"time"
)

type CreateShipRequest struct {
	ShipName      string `json:"ship_name" validate:"required"`
	Status        string `json:"status" validate:"required"`
	ShipType      string `json:"ship_type" validate:"required"`
	ShipAlias     string `json:"ship_alias" validate:"required"`
	YearOperation string `json:"year_operation" validate:"required"`
	ImageLink     string `json:"image_link" validate:"required,url"`
	Description   string `json:"description" validate:"required"`
}

type UpdateShipRequest struct {
	ID            uint   `json:"id" validate:"required"`
	ShipName      string `json:"ship_name" validate:"required"`
	Status        string `json:"status" validate:"required"`
	ShipType      string `json:"ship_type" validate:"required"`
	ShipAlias     string `json:"ship_alias" validate:"required"`
	YearOperation string `json:"year_operation" validate:"required"`
	ImageLink     string `json:"image_link" validate:"required,url"`
	Description   string `json:"description" validate:"required"`
}

type ShipResponse struct {
	ID            uint      `json:"id"`
	ShipName      string    `json:"ship_name"`
	Status        string    `json:"status"`
	ShipType      string    `json:"ship_type"`
	ShipAlias     string    `json:"ship_alias"`
	YearOperation string    `json:"year_operation"`
	ImageLink     string    `json:"image_link"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Map Ship domain to ReadShipResponse
func ShipToResponse(ship *domain.Ship) *ShipResponse {
	return &ShipResponse{
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

func ShipFromCreate(request *CreateShipRequest) *domain.Ship {
	return &domain.Ship{
		ShipName:      request.ShipName,
		ShipAlias:     request.ShipAlias,
		ShipType:      request.ShipType,
		Status:        request.Status,
		YearOperation: request.YearOperation,
		ImageLink:     request.ImageLink,
		Description:   request.Description,
	}
}

func ShipFromUpdate(request *UpdateShipRequest) *domain.Ship {
	return &domain.Ship{
		ID:            request.ID,
		ShipName:      request.ShipName,
		ShipAlias:     request.ShipAlias,
		ShipType:      request.ShipType,
		Status:        request.Status,
		YearOperation: request.YearOperation,
		ImageLink:     request.ImageLink,
		Description:   request.Description,
	}
}
