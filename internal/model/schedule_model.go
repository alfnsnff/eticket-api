package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

type ScheduleQuotaClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

// ShipDTO represents a ship.
type ScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

// QuotaDTO represents a Quota.
type ScheduleQuota struct {
	ID    uint               `json:"id"`
	Class ScheduleQuotaClass `json:"class"` // Foreign key
	Quota int                `json:"quota"`
	Price float64            `json:"price"` // Price of the quota
}

// ScheduleDTO represents a Schedule.
type ReadScheduleResponse struct {
	ID                uint            `json:"id"`
	Ship              ScheduleShip    `json:"ship"`
	DepartureHarbor   ScheduleHarbor  `json:"departure_harbor"`
	ArrivalHarbor     ScheduleHarbor  `json:"arrival_harbor"`
	DepartureDatetime time.Time       `json:"departure_datetime"`
	ArrivalDatetime   time.Time       `json:"arrival_datetime"`
	Status            string          `json:"status"` // e.g., 'active', 'inactive', 'cancelled'
	Quotas            []ScheduleQuota `json:"quotas"` // Assuming ScheduleQuota is defined elsewhere
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// ScheduleDTO represents a Schedule.
type WriteScheduleRequest struct {
	ShipID            uint      `json:"ship_id" validate:"required"`
	DepartureHarborID uint      `json:"departure_harbor_id" validate:"required"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id" validate:"required"`
	DepartureDatetime time.Time `json:"departure_datetime" validate:"required"`
	ArrivalDatetime   time.Time `json:"arrival_datetime" validate:"required,gtfield=DepartureDatetime"`
	Status            string    `json:"status" validate:"required"` // e.g., 'SCHEDULE', 'FINISHED', 'CANCELLED'
}

// ScheduleDTO represents a Schedule.
type UpdateScheduleRequest struct {
	ID                uint      `json:"id" validate:"required"`
	ShipID            uint      `json:"ship_id" validate:"required"`
	DepartureHarborID uint      `json:"departure_harbor_id" validate:"required"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id" validate:"required"`
	DepartureDatetime time.Time `json:"departure_datetime" validate:"required"`
	ArrivalDatetime   time.Time `json:"arrival_datetime" validate:"required,gtfield=DepartureDatetime"`
	Status            string    `json:"status" validate:"required"`
}

// ScheduleClassAvailability represents the availability and price for a specific class on a schedule
type ReadClassAvailabilityResponse struct {
	ClassID           uint    `json:"class_id"`
	ClassName         string  `json:"class_name"`
	Type              string  `json:"type"`
	TotalCapacity     int     `json:"total_capacity"`
	AvailableCapacity int     `json:"available_capacity"`
	Price             float64 `json:"price"`
	Currency          string  `json:"currency"` // Assuming currency is fixed or part of Fare/Route
}

// ReadScheduleDetailsWithAvailabilityResponse represents the response for schedule details with availability
type ReadScheduleDetailsResponse struct {
	ScheduleID          uint                            `json:"schedule_id"`
	DepartureHarbor     ScheduleHarbor                  `json:"departure_harbor"`
	ArrivalHarbor       ScheduleHarbor                  `json:"arrival_harbor"`
	ClassesAvailability []ReadClassAvailabilityResponse `json:"classes_availability"`
}
