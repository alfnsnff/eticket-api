package requests

import (
	"eticket-api/internal/domain"
	"time"
)

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
	CreatedAt  time.Time     `json:"created_at"`
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

func QuotaToResponse(quota *domain.Quota) *QuotaResponse {
	return &QuotaResponse{
		ID:         quota.ID,
		ScheduleID: quota.ScheduleID,
		Class: QuotaClass{
			ID:        quota.Class.ID,
			ClassName: quota.Class.ClassName,
			Type:      quota.Class.Type,
		},
		Schedule: QuotaSchedule{
			ID:                quota.Schedule.ID,
			ShipName:          quota.Schedule.Ship.ShipName,
			DepartureHarbor:   quota.Schedule.DepartureHarbor.HarborName,
			ArrivalHarbor:     quota.Schedule.ArrivalHarbor.HarborName,
			DepartureDatetime: quota.Schedule.DepartureDatetime,
			ArrivalDatetime:   quota.Schedule.ArrivalDatetime,
		},
		Price:     quota.Price,
		Quota:     quota.Quota,
		Capacity:  quota.Capacity, // Remaining quota after reservations
		CreatedAt: quota.CreatedAt,
		UpdatedAt: quota.UpdatedAt,
	}
}

func QuotaFromCreate(request *CreateQuotaRequest) *domain.Quota {
	return &domain.Quota{
		ScheduleID: request.ScheduleID,
		ClassID:    request.ClassID,
		Price:      request.Price,
		Capacity:   request.Capacity,
	}
}

func QuotaFromUpdate(request *UpdateQuotaRequest) *domain.Quota {
	return &domain.Quota{
		ID:         request.ID,
		ScheduleID: request.ScheduleID,
		ClassID:    request.ClassID,
		Price:      request.Price,
		Capacity:   request.Capacity,
	}
}
