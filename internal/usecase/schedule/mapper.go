package schedule

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Schedule entity to ReadScheduleResponse model
func ToReadScheduleResponse(schedule *entity.Schedule) *model.ReadScheduleResponse {
	if schedule == nil {
		return nil
	}

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

// Map a slice of Schedule entities to ReadScheduleResponse models
func ToReadScheduleResponses(schedules []*entity.Schedule) []*model.ReadScheduleResponse {
	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = ToReadScheduleResponse(schedule)
	}
	return responses
}

// Map WriteScheduleRequest model to Schedule entity
func FromWriteScheduleRequest(request *model.WriteScheduleRequest) *entity.Schedule {
	return &entity.Schedule{
		RouteID:           request.RouteID,
		ShipID:            request.ShipID,
		DepartureDatetime: &request.DepartureDatetime,
		ArrivalDatetime:   &request.ArrivalDatetime,
		Status:            &request.Status,
	}
}

// Map UpdateScheduleRequest model to Schedule entity
func FromUpdateScheduleRequest(request *model.UpdateScheduleRequest, schedule *entity.Schedule) {
	schedule.RouteID = request.RouteID
	schedule.ShipID = request.ShipID
	schedule.DepartureDatetime = &request.DepartureDatetime
	schedule.ArrivalDatetime = &request.ArrivalDatetime
	schedule.Status = &request.Status
}
