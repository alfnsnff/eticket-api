package route

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Route entity to ReadRouteResponse model
func ToReadRouteResponse(route *entity.Route) *model.ReadRouteResponse {
	if route == nil {
		return nil
	}

	return &model.ReadRouteResponse{
		ID: route.ID,
		DepartureHarbor: model.RouteHarbor{
			ID:         route.DepartureHarbor.ID,
			HarborName: route.DepartureHarbor.HarborName,
		},
		ArrivalHarbor: model.RouteHarbor{
			ID:         route.ArrivalHarbor.ID,
			HarborName: route.ArrivalHarbor.HarborName,
		},
		CreatedAt: route.CreatedAt,
		UpdatedAt: route.UpdatedAt,
	}
}

// Map a slice of Route entities to ReadRouteResponse models
func ToReadRouteResponses(routes []*entity.Route) []*model.ReadRouteResponse {
	responses := make([]*model.ReadRouteResponse, len(routes))
	for i, route := range routes {
		responses[i] = ToReadRouteResponse(route)
	}
	return responses
}

// Map WriteRouteRequest model to Route entity
func FromWriteRouteRequest(request *model.WriteRouteRequest) *entity.Route {
	return &entity.Route{
		DepartureHarborID: request.DepartureHarborID,
		ArrivalHarborID:   request.ArrivalHarborID,
	}
}

// Map UpdateRouteRequest model to Route entity
func FromUpdateRouteRequest(request *model.UpdateRouteRequest, route *entity.Route) {
	route.DepartureHarborID = request.DepartureHarborID
	route.ArrivalHarborID = request.ArrivalHarborID
}
