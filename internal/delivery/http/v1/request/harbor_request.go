package request

import (
	"time"
)

type CreateHarborRequest struct {
	HarborName    string `json:"harbor_name" validate:"required"`
	Status        string `json:"status" validate:"required"`
	YearOperation string `json:"year_operation" validate:"required"`
	HarborAlias   string `json:"harbor_alias" validate:"required"`
}

type UpdateHarborRequest struct {
	ID            uint   `json:"id" validate:"required"`
	HarborName    string `json:"harbor_name" validate:"required"`
	Status        string `json:"status" validate:"required"`
	YearOperation string `json:"year_operation" validate:"required"`
	HarborAlias   string `json:"harbor_alias" validate:"required"`
}

type HarborResponse struct {
	ID            uint      `json:"id"`
	HarborName    string    `json:"harbor_name"`
	Status        string    `json:"status"`
	HarborAlias   string    `json:"harbor_alias"`
	YearOperation string    `json:"year_operation"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
