package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToScheduleModel(schedule *entities.Schedule) *model.ReadScheduleResponse {
	response := new(model.ReadScheduleResponse)
	copier.Copy(&response, &schedule) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToSchedulesModel(schedules []*entities.Schedule) []*model.ReadScheduleResponse {
	responses := []*model.ReadScheduleResponse{}
	for _, schedule := range schedules {
		responses = append(responses, ToScheduleModel(schedule))
	}
	return responses
}

func ToScheduleEntity(request *model.WriteScheduleRequest) *entities.Schedule {
	schedule := new(entities.Schedule)
	copier.Copy(&schedule, &request)
	return schedule
}
