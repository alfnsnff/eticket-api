package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	DB                 *gorm.DB
	ScheduleRepository *repository.ScheduleRepository
	RouteRepository    *repository.RouteRepository
	FareRepository     *repository.FareRepository
	TicketRepository   *repository.TicketRepository
}

func NewScheduleUsecase(
	db *gorm.DB,
	schedule_repository *repository.ScheduleRepository,
	route_repository *repository.RouteRepository,
	fare_repository *repository.FareRepository,
	ticket_repository *repository.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		DB:                 db,
		ScheduleRepository: schedule_repository,
		RouteRepository:    route_repository,
		FareRepository:     fare_repository,
		TicketRepository:   ticket_repository,
	}
}

func (sc *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.WriteScheduleRequest) error {
	schedule := mapper.ScheduleMapper.FromWrite(request)

	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		return sc.ScheduleRepository.Create(tx, schedule)
	})
}

func (sc *ScheduleUsecase) GetAllSchedules(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	schedules := []*entity.Schedule{}

	err := tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		var err error
		schedules, err = sc.ScheduleRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	return mapper.ScheduleMapper.ToModels(schedules), nil
}

func (sc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	schedule := new(entity.Schedule)

	err := tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		var err error
		schedule, err = sc.ScheduleRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}

	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	return mapper.ScheduleMapper.ToModel(schedule), nil
}

func (sc *ScheduleUsecase) UpdateSchedule(ctx context.Context, id uint, request *model.UpdateScheduleRequest) error {
	schedule := mapper.ScheduleMapper.FromUpdate(request)
	schedule.ID = id

	if schedule.ID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}

	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		return sc.ScheduleRepository.Update(tx, schedule)
	})
}

func (sc *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {

	return tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		schedule, err := sc.ScheduleRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}
		return sc.ScheduleRepository.Delete(tx, schedule)
	})

}

// // SearchSchedule searches a schedule by departure, arrival, and date
// func (s *ScheduleUsecase) SearchSchedule(ctx context.Context, request *model.ScheduleSearchRequest) (*entity.Schedule, error) {
// 	var schedule *entity.Schedule

// 	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
// 		schedule, err = s.ScheduleRepository.Search(txDB, route.ID, req.Date, req.ShipID)
// 		return err
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search schedule: %w", err)
// 	}

// 	return schedule, nil
// }
