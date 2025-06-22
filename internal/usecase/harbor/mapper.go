package harbor

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Harbor entity to ReadHarborResponse model
func ToReadHarborResponse(harbor *entity.Harbor) *model.ReadHarborResponse {
	if harbor == nil {
		return nil
	}

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

// Map a slice of Harbor entities to ReadHarborResponse models
func ToReadHarborResponses(harbors []*entity.Harbor) []*model.ReadHarborResponse {
	responses := make([]*model.ReadHarborResponse, len(harbors))
	for i, harbor := range harbors {
		responses[i] = ToReadHarborResponse(harbor)
	}
	return responses
}

// Map WriteHarborRequest model to Harbor entity
func FromWriteHarborRequest(request *model.WriteHarborRequest) *entity.Harbor {
	return &entity.Harbor{
		HarborName:    request.HarborName,
		Status:        request.Status,
		HarborAlias:   request.HarborAlias,
		YearOperation: request.YearOperation,
	}
}

// Map UpdateHarborRequest model to Harbor entity
func FromUpdateHarborRequest(request *model.UpdateHarborRequest, harbor *entity.Harbor) {
	harbor.HarborName = request.HarborName
	harbor.Status = request.Status
	harbor.HarborAlias = request.HarborAlias
	harbor.YearOperation = request.YearOperation
}
