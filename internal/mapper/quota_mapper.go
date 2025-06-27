package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

func QuotaToResponse(quota *domain.Quota) *model.ReadQuotaResponse {
	return &model.ReadQuotaResponse{
		ID:         quota.ID,
		ScheduleID: quota.ScheduleID,
		Class: model.QuotaClass{
			ID:        quota.Class.ID,
			ClassName: quota.Class.ClassName,
			Type:      quota.Class.Type,
		},
		Schedule: model.QuotaSchedule{
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
