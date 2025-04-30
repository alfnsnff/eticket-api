package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToRouteModel(route *entities.Route) *model.ReadRouteResponse {
	response := new(model.ReadRouteResponse)
	copier.Copy(&response, &route) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToRoutesModel(routes []*entities.Route) []*model.ReadRouteResponse {
	responses := []*model.ReadRouteResponse{}
	for _, route := range routes {
		responses = append(responses, ToRouteModel(route))
	}
	return responses
}

func ToRouteEntity(request *model.WriteRouteRequest) *entities.Route {
	route := new(entities.Route)
	copier.Copy(&route, &request)
	return route
}
