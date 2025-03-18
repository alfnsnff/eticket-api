package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"fmt"
)

type ScheduleUsecase struct {
	ScheduleRepository entities.ScheduleRepositoryInterface
}

func NewScheduleUsecase(scheduleRepository entities.ScheduleRepositoryInterface) ScheduleUsecase {
	return ScheduleUsecase{ScheduleRepository: scheduleRepository}
}

// Createschedule validates and creates a new schedule
func (s *ScheduleUsecase) CreateSchedule(schedule *entities.Schedule) error {
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}
	return s.ScheduleRepository.Create(schedule)
}

// GetAllschedulees retrieves all schedules
func (s *ScheduleUsecase) GetAllSchedules() ([]*entities.Schedule, error) {
	return s.ScheduleRepository.GetAll()
}

// GetscheduleByID retrieves a schedule by its ID
func (s *ScheduleUsecase) GetScheduleByID(id uint) (*entities.Schedule, error) {
	schedule, err := s.ScheduleRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("schedule not found")
	}
	return schedule, nil
}

// Updateschedule updates an existing schedule
func (s *ScheduleUsecase) UpdateSchedule(schedule *entities.Schedule) error {
	if schedule.ID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}
	return s.ScheduleRepository.Update(schedule)
}

// Deleteschedule deletes a schedule by its ID
func (s *ScheduleUsecase) DeleteSchedule(id uint) error {
	schedule, err := s.ScheduleRepository.GetByID(id)
	if err != nil {
		return err
	}
	if schedule == nil {
		return errors.New("schedule not found")
	}
	return s.ScheduleRepository.Delete(id)
}
