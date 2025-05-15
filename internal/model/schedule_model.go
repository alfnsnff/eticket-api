package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type ScheduleRoute struct {
	ID              uint           `json:"id"`
	DepartureHarbor ScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor   ScheduleHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type ScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// ScheduleDTO represents a Schedule.
type ReadScheduleResponse struct {
	ID               uint          `json:"id"`
	Ship             ScheduleShip  `json:"ship"`
	Route            ScheduleRoute `json:"route"`
	ScheduleDatetime time.Time     `json:"departure_datetime"`
	ArrivalDatetime  time.Time     `json:"arrival_datetime"`
	Status           string        `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// ScheduleDTO represents a Schedule.
type WriteScheduleRequest struct {
	RouteID          uint      `json:"route_id"`
	ShipID           uint      `json:"ship_id"`
	ScheduleDatetime time.Time `json:"departure_datetime"`
	ArrivalDatetime  time.Time `json:"arrival_datetime"`
	Status           string    `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
}

// ScheduleDTO represents a Schedule.
type UpdateScheduleRequest struct {
	ID               uint      `json:"id"`
	RouteID          uint      `json:"route_id"`
	ShipID           uint      `json:"ship_id"`
	ScheduleDatetime time.Time `json:"departure_datetime"`
	ArrivalDatetime  time.Time `json:"arrival_datetime"`
	Status           string    `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
}

// ScheduleClassAvailability represents the availability and price for a specific class on a schedule
type ScheduleClassAvailability struct {
	ClassID           uint    `json:"class_id"`
	ClassName         string  `json:"class_name"`
	TotalCapacity     int     `json:"total_capacity"`
	AvailableCapacity int     `json:"available_capacity"`
	Price             float32 `json:"price"`
	Currency          string  `json:"currency"` // Assuming currency is fixed or part of Fare/Route
}

// ReadScheduleDetailsWithAvailabilityResponse represents the response for schedule details with availability
type ReadScheduleDetailsWithAvailabilityResponse struct {
	ScheduleID          uint                        `json:"schedule_id"`
	RouteID             uint                        `json:"route_id"`
	ClassesAvailability []ScheduleClassAvailability `json:"classes_availability"`
}
