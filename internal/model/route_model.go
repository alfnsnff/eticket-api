package model

import (
	"time"
)

// HarborDTO represents a harbor.
type RouteHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type ReadRouteResponse struct {
	ID              uint           `json:"id"`
	DepartureHarbor ScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   ScheduleHarbor `json:"arrival_harbor"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
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
