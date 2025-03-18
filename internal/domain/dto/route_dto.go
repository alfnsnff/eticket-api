package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type RouteHarborRes struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type RouteRes struct {
	ID              uint            `json:"id"`
	DepartureHarbor TicketHarborRes `json:"departure_harbor"`
	ArrivalHarbor   TicketHarborRes `json:"arrival_harbor"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func ToRouteDTO(route *entities.Route) RouteRes {
	var routeResponse RouteRes
	copier.Copy(&routeResponse, &route) // Automatically maps matching fields
	return routeResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToRouteDTOs(routes []*entities.Route) []RouteRes {
	var routeResponses []RouteRes
	for _, route := range routes {
		routeResponses = append(routeResponses, ToRouteDTO(route))
	}
	return routeResponses
}
