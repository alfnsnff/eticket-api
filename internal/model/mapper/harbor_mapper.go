package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToHarborModel(harbor *entities.Harbor) *model.ReadHarborResponse {
	response := new(model.ReadHarborResponse)
	copier.Copy(&response, &harbor) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToHarborsModel(harbors []*entities.Harbor) []*model.ReadHarborResponse {
	responses := []*model.ReadHarborResponse{}
	for _, harbor := range harbors {
		responses = append(responses, ToHarborModel(harbor))
	}
	return responses
}

func ToHarborEntity(request *model.WriteHarborRequest) *entities.Harbor {
	harbor := new(entities.Harbor)
	copier.Copy(&harbor, &request)
	return harbor
}
