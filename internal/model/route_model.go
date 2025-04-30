package model

import (
	"time"
)

// HarborDTO represents a harbor.
type RouteHarbor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RouteDTO represents a travel route.
type ReadRouteResponse struct {
	ID              uint        `json:"id"`
	DepartureHarbor RouteHarbor `json:"departure_harbor"`
	ArrivalHarbor   RouteHarbor `json:"arrival_harbor"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type WriteRouteRequest struct {
	ID                uint `json:"id"`
	DepartureHarborID uint `json:"departure_harbor_id"`
	ArrivalHarborID   uint `json:"arrival_harbor_id"`
}
