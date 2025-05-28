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
	ID    uint              `json:"id"`
	Class FareManifestClass `json:"class"`
	Ship  FareManifestShip  `json:"ship"`
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
	ID              uint           `json:"id"`
	DepartureHarbor ScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   ScheduleHarbor `json:"arrival_harbor"`
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
	RouteID     uint    `json:"route_id"`
	ManifestID  uint    `json:"manifest_id"`
	TicketPrice float32 `json:"ticket_price"`
}

type UpdateFareRequest struct {
	ID          uint    `json:"id"`
	RouteID     uint    `json:"route_id"`
	ManifestID  uint    `json:"manifest_id"`
	TicketPrice float32 `json:"ticket_price"`
}
