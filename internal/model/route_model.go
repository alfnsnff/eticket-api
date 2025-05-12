package model

import (
	"time"
)

// RouteDTO represents a travel route.
type ReadRouteResponse struct {
	ID                uint      `json:"id"`
	DepartureHarborID uint      `json:"departure_harbor_id"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type WriteRouteRequest struct {
	DepartureHarborID uint `json:"departure_harbor_id"`
	ArrivalHarborID   uint `json:"arrival_harbor_id"`
}

type UpdateRouteRequest struct {
	ID                uint `json:"id"`
	DepartureHarborID uint `json:"departure_harbor_id"`
	ArrivalHarborID   uint `json:"arrival_harbor_id"`
}
