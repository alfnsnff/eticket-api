package request

import (
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

type ReadShipResponse struct {
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
