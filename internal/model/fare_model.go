package model

import (
	"time"
)

type FareManifestClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
}

type FareManifestShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

type FareManifest struct {
	ID       uint              `json:"id"`
	Class    FareManifestClass `json:"class"`
	Ship     FareManifestShip  `json:"ship"`
	Capacity int               `json:"capacity"`
}

type CountFareRequest struct {
	ID         uint `json:"id"`
	RouteID    uint `json:"route_id"`
	ManifestID uint `json:"manifest_id"`
}

// HarborDTO represents a harbor.
type FareRouteHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type FareRoute struct {
	ID              uint            `json:"id"`
	DepartureHarbor FareRouteHarbor `json:"departure_harbor"`
	ArrivalHarbor   FareRouteHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a Ship.
type ReadFareResponse struct {
	ID          uint         `json:"id"`
	Route       FareRoute    `json:"route"`
	Manifest    FareManifest `json:"manifest"`
	TicketPrice float32      `json:"ticket_price"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type WriteFareRequest struct {
	RouteID     uint    `json:"route_id" validate:"required"`
	ManifestID  uint    `json:"manifest_id" validate:"required"`
	TicketPrice float32 `json:"ticket_price" validate:"required,gt=0"`
}

type UpdateFareRequest struct {
	ID          uint    `json:"id" validate:"required"`
	RouteID     uint    `json:"route_id" validate:"required"`
	ManifestID  uint    `json:"manifest_id" validate:"required"`
	TicketPrice float32 `json:"ticket_price" validate:"required,gt=0"`
}
