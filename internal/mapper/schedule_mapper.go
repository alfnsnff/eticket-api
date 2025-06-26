package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Schedule domain to ReadScheduleResponse model
func ScheduleToResponse(schedule *domain.Schedule) *model.ReadScheduleResponse {

	quotas := make([]model.ScheduleQuota, len(schedule.Quotas))
	for i, quota := range schedule.Quotas {
		quotas[i] = model.ScheduleQuota{
			ID: quota.ID,
			Class: model.ScheduleQuotaClass{
				ID:        quota.Class.ID,
				ClassName: quota.Class.ClassName,
				Type:      quota.Class.Type,
			},
			Quota: quota.Quota,
			Price: quota.Price,
		}
	}
	return &model.ReadScheduleResponse{
		ID: schedule.ID,
		Ship: model.ScheduleShip{
			ID:       schedule.Ship.ID,
			ShipName: schedule.Ship.ShipName,
		},
		DepartureHarbor: model.ScheduleHarbor{
			ID:         schedule.DepartureHarbor.ID,
			HarborName: schedule.DepartureHarbor.HarborName,
		},
		ArrivalHarbor: model.ScheduleHarbor{
			ID:         schedule.ArrivalHarbor.ID,
			HarborName: schedule.ArrivalHarbor.HarborName,
		},
		Quotas:            quotas,
		DepartureDatetime: schedule.DepartureDatetime,
		ArrivalDatetime:   schedule.ArrivalDatetime,
		Status:            schedule.Status,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
	}
}
