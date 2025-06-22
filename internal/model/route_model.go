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
	ID              uint        `json:"id"`
	DepartureHarbor RouteHarbor `json:"departure_harbor"`
	ArrivalHarbor   RouteHarbor `json:"arrival_harbor"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type WriteRouteRequest struct {
	DepartureHarborID uint `json:"departure_harbor_id" validate:"required"`
	ArrivalHarborID   uint `json:"arrival_harbor_id" validate:"required,nefield=DepartureHarborID"`
}

type UpdateRouteRequest struct {
	ID                uint `json:"id" validate:"required"`
	DepartureHarborID uint `json:"departure_harbor_id" validate:"required"`
	ArrivalHarborID   uint `json:"arrival_harbor_id" validate:"required,nefield=DepartureHarborID"`
}
