package model

import (
	"time"
)

// // HarborDTO represents a harbor.
// type ScheduleHarbor struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // RouteDTO represents a travel route.
// type ScheduleRoute struct {
// 	ID              uint           `json:"id"`
// 	DepartureHarbor ScheduleHarbor `json:"departure_harbor"`
// 	ArrivalHarbor   ScheduleHarbor `json:"arrival_harbor"`
// }

// // ShipDTO represents a ship.
// type ScheduleShip struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// type SearchScheduleRequest struct {
// 	DepartureHarborID uint      `json:"departure_harbor_id"`
// 	ArrivalHarborID   uint      `json:"arrival_harbor_id"`
// 	Date              time.Time `json:"date"`
// 	ShipID            *uint     `json:"ship_id,omitempty"` // optional
// }

// type ScheduleQuotaResponse struct {
// 	FareID    uint    `json:"fare_id"`
// 	ClassName string  `json:"class_name"`
// 	Price     float32 `json:"price"`
// 	Capacity  int     `json:"capacity"`
// 	Booked    int     `json:"booked"`
// 	Available int     `json:"available"`
// }

// ScheduleDTO represents a Schedule.
type ReadScheduleResponse struct {
	ID        uint      `json:"id"`
	ShipID    uint      `json:"ship"`
	RouteID   uint      `json:"route"`
	Datetime  time.Time `json:"datetime"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScheduleDTO represents a Schedule.
type WriteScheduleRequest struct {
	RouteID  uint      `json:"route_id"`
	ShipID   uint      `json:"ship_id"`
	Datetime time.Time `json:"datetime"`
}

// ScheduleDTO represents a Schedule.
type UpdateScheduleRequest struct {
	ID       uint      `json:"id"`
	RouteID  uint      `json:"route_id"`
	ShipID   uint      `json:"ship_id"`
	Datetime time.Time `json:"datetime"`
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
	ShipName            string                      `json:"ship_name"`
	ClassesAvailability []ScheduleClassAvailability `json:"classes_availability"`
}
