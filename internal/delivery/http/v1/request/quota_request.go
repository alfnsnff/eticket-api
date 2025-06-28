package request

import "time"

type CreateQuotaRequest struct {
	ScheduleID uint    `json:"schedule_id" validate:"required,gt=0"`
	ClassID    uint    `json:"class_id" validate:"required,gt=0"`
	Capacity   int     `json:"capacity" validate:"required,gte=0"`
	Price      float64 `json:"price"`
}

type UpdateQuotaRequest struct {
	ID         uint    `json:"id" validate:"required,gt=0"`
	ScheduleID uint    `json:"schedule_id" validate:"required,gt=0"`
	ClassID    uint    `json:"class_id" validate:"required,gt=0"`
	Price      float64 `json:"price"`
	Capacity   int     `json:"capacity" validate:"required,gte=0"`
}

type QuotaResponse struct {
	ID         uint          `json:"id"`
	ScheduleID uint          `json:"schedule_id"`
	Class      QuotaClass    `json:"class"`
	Schedule   QuotaSchedule `json:"schedule"`
	Price      float64       `json:"price"`
	Quota      int           `json:"quota"`
	Capacity   int           `json:"Capacity"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type QuotaClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type QuotaSchedule struct {
	ID                uint      `json:"id"`
	ShipName          string    `json:"ship_name"`
	DepartureHarbor   string    `json:"departure_harbor"`
	ArrivalHarbor     string    `json:"arrival_harbor"`
	DepartureDatetime time.Time `json:"departure_datetime"`
	ArrivalDatetime   time.Time `json:"arrival_datetime"`
}
