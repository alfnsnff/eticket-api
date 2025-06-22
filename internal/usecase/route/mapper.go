package route

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Route domain to ReadRouteResponse model
func RouteToResponse(route *domain.Route) *model.ReadRouteResponse {
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
