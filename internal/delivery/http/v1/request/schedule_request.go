package request

import (
	"time"
)

type CreateScheduleRequest struct {
	ShipID            uint      `json:"ship_id" validate:"required"`
	DepartureHarborID uint      `json:"departure_harbor_id" validate:"required"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id" validate:"required"`
	DepartureDatetime time.Time `json:"departure_datetime" validate:"required"`
	ArrivalDatetime   time.Time `json:"arrival_datetime" validate:"required,gtfield=DepartureDatetime"`
	Status            string    `json:"status" validate:"required"`
}

type UpdateScheduleRequest struct {
	ID                uint      `json:"id" validate:"required"`
	ShipID            uint      `json:"ship_id" validate:"required"`
	DepartureHarborID uint      `json:"departure_harbor_id" validate:"required"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id" validate:"required"`
	DepartureDatetime time.Time `json:"departure_datetime" validate:"required"`
	ArrivalDatetime   time.Time `json:"arrival_datetime" validate:"required,gtfield=DepartureDatetime"`
	Status            string    `json:"status" validate:"required"`
}

type ScheduleResponse struct {
	ID                uint            `json:"id"`
	Ship              ScheduleShip    `json:"ship"`
	DepartureHarbor   ScheduleHarbor  `json:"departure_harbor"`
	ArrivalHarbor     ScheduleHarbor  `json:"arrival_harbor"`
	DepartureDatetime time.Time       `json:"departure_datetime"`
	ArrivalDatetime   time.Time       `json:"arrival_datetime"`
	Status            string          `json:"status"`
	Quotas            []ScheduleQuota `json:"quotas"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type ScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

type ScheduleQuotaClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type ScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

type ScheduleQuota struct {
	ID    uint               `json:"id"`
	Class ScheduleQuotaClass `json:"class"`
	Quota int                `json:"quota"`
	Price float64            `json:"price"`
}
