package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToClassModel(class *entities.Class) *model.ReadClassResponse {
	response := new(model.ReadClassResponse)
	copier.Copy(&response, &class) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToClassesModel(classes []*entities.Class) []*model.ReadClassResponse {
	responses := []*model.ReadClassResponse{}
	for _, class := range classes {
		responses = append(responses, ToClassModel(class))
	}
	return responses
}

func ToClassEntity(request *model.WriteClassRequest) *entities.Class {
	class := new(entities.Class)
	copier.Copy(&class, &request)
	return class
}
