package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ScheduleHarbor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RouteDTO represents a travel route.
type ScheduleRoute struct {
	ID              uint           `json:"id"`
	DepartureHarbor ScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   ScheduleHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type ScheduleShip struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ScheduleDTO represents a Schedule.
type ReadScheduleResponse struct {
	ID        uint          `json:"id"`
	Datetime  time.Time     `json:"datetime"`
	Ship      ScheduleShip  `json:"ship"`
	Route     ScheduleRoute `json:"route"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ScheduleDTO represents a Schedule.
type WriteScheduleRequest struct {
	ID       uint      `json:"id"`
	RouteID  uint      `json:"route_id"`
	ShipID   uint      `json:"ship_id"`
	Datetime time.Time `json:"datetime"`
}

type SearchScheduleRequest struct {
	DepartureHarborID uint      `json:"departure_harbor_id"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id"`
	Date              time.Time `json:"date"`
	ShipID            *uint     `json:"ship_id,omitempty"` // optional
}

type ScheduleQuotaResponse struct {
	PriceID   uint    `json:"price_id"`
	ClassName string  `json:"class_name"`
	Price     float32 `json:"price"`
	Capacity  int     `json:"capacity"`
	Booked    int     `json:"booked"`
	Available int     `json:"available"`
}
