package requests

import (
	"eticket-api/internal/domain"
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

// Helper to build quotas slice
func buildScheduleQuotas(quotasDomain []*domain.Quota) []ScheduleQuota {
	quotas := make([]ScheduleQuota, len(quotasDomain))
	for i, quota := range quotasDomain {
		quotas[i] = ScheduleQuota{
			ID: quota.ID,
			Class: ScheduleQuotaClass{
				ID:        quota.Class.ID,
				ClassName: quota.Class.ClassName,
				Type:      quota.Class.Type,
			},
			Quota: quota.Quota,
			Price: quota.Price,
		}
	}
	return quotas
}

// Map Schedule domain to ReadScheduleResponse model
func ScheduleToResponse(schedule *domain.Schedule) *ScheduleResponse {
	return &ScheduleResponse{
		ID: schedule.ID,
		Ship: ScheduleShip{
			ID:       schedule.Ship.ID,
			ShipName: schedule.Ship.ShipName,
		},
		DepartureHarbor: ScheduleHarbor{
			ID:         schedule.DepartureHarbor.ID,
			HarborName: schedule.DepartureHarbor.HarborName,
		},
		ArrivalHarbor: ScheduleHarbor{
			ID:         schedule.ArrivalHarbor.ID,
			HarborName: schedule.ArrivalHarbor.HarborName,
		},
		Quotas:            buildScheduleQuotas(schedule.Quotas),
		DepartureDatetime: schedule.DepartureDatetime,
		ArrivalDatetime:   schedule.ArrivalDatetime,
		Status:            schedule.Status,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
	}
}

func ScheduleFromCreate(request *CreateScheduleRequest) *domain.Schedule {
	return &domain.Schedule{
		ShipID:            request.ShipID,
		DepartureHarborID: request.DepartureHarborID,
		ArrivalHarborID:   request.ArrivalHarborID,
		DepartureDatetime: request.DepartureDatetime,
		ArrivalDatetime:   request.ArrivalDatetime,
		Status:            request.Status,
	}
}

func ScheduleFromUpdate(request *UpdateScheduleRequest) *domain.Schedule {
	return &domain.Schedule{
		ID:                request.ID,
		ShipID:            request.ShipID,
		DepartureHarborID: request.DepartureHarborID,
		ArrivalHarborID:   request.ArrivalHarborID,
		DepartureDatetime: request.DepartureDatetime,
		ArrivalDatetime:   request.ArrivalDatetime,
		Status:            request.Status,
	}
}
