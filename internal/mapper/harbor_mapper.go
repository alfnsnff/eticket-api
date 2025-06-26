package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Harbor domain to ReadHarborResponse model
func HarborToResponse(harbor *domain.Harbor) *model.ReadHarborResponse {
	return &model.ReadHarborResponse{
		ID:            harbor.ID,
		HarborName:    harbor.HarborName,
		Status:        harbor.Status,
		HarborAlias:   harbor.HarborAlias,
		YearOperation: harbor.YearOperation,
		CreatedAt:     harbor.CreatedAt,
		UpdatedAt:     harbor.UpdatedAt,
	}
}
