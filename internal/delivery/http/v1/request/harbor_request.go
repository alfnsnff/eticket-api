package requests

import (
	"eticket-api/internal/domain"
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

// Map Harbor domain to ReadHarborResponse model
func HarborToResponse(harbor *domain.Harbor) *HarborResponse {
	return &HarborResponse{
		ID:            harbor.ID,
		HarborName:    harbor.HarborName,
		Status:        harbor.Status,
		HarborAlias:   harbor.HarborAlias,
		YearOperation: harbor.YearOperation,
		CreatedAt:     harbor.CreatedAt,
		UpdatedAt:     harbor.UpdatedAt,
	}
}

func HarborFromCreate(request *CreateHarborRequest) *domain.Harbor {
	return &domain.Harbor{
		HarborName:    request.HarborName,
		HarborAlias:   request.HarborAlias,
		Status:        request.Status,
		YearOperation: request.YearOperation,
	}
}

func HarborFromUpdate(request *UpdateHarborRequest) *domain.Harbor {
	return &domain.Harbor{
		ID:            request.ID,
		HarborName:    request.HarborName,
		HarborAlias:   request.HarborAlias,
		Status:        request.Status,
		YearOperation: request.YearOperation,
	}
}
