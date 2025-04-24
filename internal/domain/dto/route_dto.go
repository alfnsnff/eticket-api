package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type RouteHarbor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RouteDTO represents a travel route.
type RouteRead struct {
	ID              uint        `json:"id"`
	DepartureHarbor RouteHarbor `json:"departure_harbor"`
	ArrivalHarbor   RouteHarbor `json:"arrival_harbor"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type RouteCreate struct {
	DepartureHarborID uint `json:"departure_harbor_id"`
	ArrivalHarborID   uint `json:"arrival_harbor_id"`
}

func ToRouteDTO(route *entities.Route) RouteRead {
	var routeResponse RouteRead
	copier.Copy(&routeResponse, &route) // Automatically maps matching fields
	return routeResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToRouteDTOs(routes []*entities.Route) []RouteRead {
	var routeResponses []RouteRead
	for _, route := range routes {
		routeResponses = append(routeResponses, ToRouteDTO(route))
	}
	return routeResponses
}

func ToRouteEntity(routeCreate *RouteCreate) entities.Route {
	var route entities.Route
	copier.Copy(&route, &routeCreate)
	return route
}
