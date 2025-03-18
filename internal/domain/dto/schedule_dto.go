package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type ScheduleHarborRes struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type ScheduleRouteRes struct {
	ID              uint              `json:"id"`
	DepartureHarbor ScheduleHarborRes `json:"departure_harbor"`
	ArrivalHarbor   ScheduleHarborRes `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type ScheduleShipRes struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
	Capacity uint   `json:"capacity"`
}

// ScheduleDTO represents a Schedule.
type ScheduleRes struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Datetime  time.Time        `gorm:"not null" json:"datetime"`
	Ship      ScheduleShipRes  `json:"ship"`
	Route     ScheduleRouteRes `json:"route"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func ToScheduleDTO(schedule *entities.Schedule) ScheduleRes {
	var scheduleResponse ScheduleRes
	copier.Copy(&scheduleResponse, &schedule) // Automatically maps matching fields
	return scheduleResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToScheduleDTOs(schedules []*entities.Schedule) []ScheduleRes {
	var scheduleResponses []ScheduleRes
	for _, schedule := range schedules {
		scheduleResponses = append(scheduleResponses, ToScheduleDTO(schedule))
	}
	return scheduleResponses
}
