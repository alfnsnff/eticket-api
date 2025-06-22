package schedule

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Schedule domain to ReadScheduleResponse model
func ScheduleToResponse(schedule *domain.Schedule) *model.ReadScheduleResponse {
	return &model.ReadScheduleResponse{
		ID: schedule.ID,
		Ship: model.ScheduleShip{
			ID:       schedule.Ship.ID,
			ShipName: schedule.Ship.ShipName,
		},
		Route: model.ScheduleRoute{
			ID: schedule.Route.ID,
			DepartureHarbor: model.ScheduleHarbor{
				ID:         schedule.Route.DepartureHarbor.ID,
				HarborName: schedule.Route.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.ScheduleHarbor{
				ID:         schedule.Route.ArrivalHarbor.ID,
				HarborName: schedule.Route.ArrivalHarbor.HarborName,
			},
		},
		DepartureDatetime: *schedule.DepartureDatetime,
		ArrivalDatetime:   *schedule.ArrivalDatetime,
		Status:            *schedule.Status,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
	}
}
