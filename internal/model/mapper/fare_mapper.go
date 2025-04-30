package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToFareModel(fare *entities.Fare) *model.ReadFareResponse {
	response := new(model.ReadFareResponse)
	copier.Copy(&response, &fare) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToFaresModel(fares []*entities.Fare) []*model.ReadFareResponse {
	responses := []*model.ReadFareResponse{}
	for _, fare := range fares {
		responses = append(responses, ToFareModel(fare))
	}
	return responses
}

func ToFareEntity(request *model.WriteFareRequest) *entities.Fare {
	fare := new(entities.Fare)
	copier.Copy(&fare, &request)
	return fare
}
