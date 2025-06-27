package model

import "time"

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

// QuotaDTO represents a Quota.
type ReadQuotaResponse struct {
	ID             uint          `json:"id"`
	ScheduleID     uint          `json:"schedule_id"` // Foreign key
	Class          QuotaClass    `json:"class"`       // Foreign key
	Schedule       QuotaSchedule `json:"schedule"`    // Foreign key
	Price          float64       `json:"price"`       // Price of the quota
	Quota          int           `json:"quota"`
	RemainingQuota int           `json:"remaining_quota"` // Remaining quota after reservations
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// QuotaDTO represents a Quota.
type WriteQuotaRequest struct {
	ScheduleID uint    `json:"schedule_id" validate:"required,gt=0"`
	ClassID    uint    `json:"class_id" validate:"required,gt=0"` // must be > 0
	Price      float64 `json:"price"`                             // Price of the quota
	Quota      int     `json:"quota" validate:"required,gte=0"`   // must be ≥ 0
}

// QuotaDTO represents a Quota.
type UpdateQuotaRequest struct {
	ID         uint    `json:"id" validate:"required,gt=0"`          // must be > 0
	ScheduleID uint    `json:"schedule_id" validate:"required,gt=0"` // must be > 0
	ClassID    uint    `json:"class_id" validate:"required,gt=0"`    // must be > 0
	Price      float64 `json:"price"`                                // Price of the quota
	Quota      int     `json:"quota" validate:"required,gte=0"`      // must be ≥ 0
}
