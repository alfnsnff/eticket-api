package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	DB                 *gorm.DB // Assuming you have a DB field for the transaction manager
	ClassRepository    domain.ClassRepository
	ShipRepository     domain.ShipRepository
	ScheduleRepository domain.ScheduleRepository
	TicketRepository   domain.TicketRepository
}

func NewScheduleUsecase(
	db *gorm.DB,
	class_repository domain.ClassRepository,
	ship_repository domain.ShipRepository,
	schedule_repository domain.ScheduleRepository,
	ticket_repository domain.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		DB:                 db,
		ClassRepository:    class_repository,
		ShipRepository:     ship_repository,
		ScheduleRepository: schedule_repository,
		TicketRepository:   ticket_repository,
	}
}

func (sc *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.WriteScheduleRequest) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule := &domain.Schedule{
		ShipID:            request.ShipID,
		DepartureHarborID: request.DepartureHarborID,
		ArrivalHarborID:   request.ArrivalHarborID,
		DepartureDatetime: request.DepartureDatetime,
		ArrivalDatetime:   request.ArrivalDatetime,
		Status:            request.Status,
	}
	if err := sc.ScheduleRepository.Insert(tx, schedule); err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sc *ScheduleUsecase) ListSchedules(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadScheduleResponse, int, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()
	total, err := sc.ScheduleRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	schedules, err := sc.ScheduleRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all schedules: %w", err)
	}

	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = mapper.ScheduleToResponse(schedule)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (sc *ScheduleUsecase) ListActiveSchedules(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	schedules, err := sc.ScheduleRepository.FindActiveSchedule(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all active schedules: %w", err)
	}

	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = mapper.ScheduleToResponse(schedule)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, nil
}

func (sc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule, err := sc.ScheduleRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ScheduleToResponse(schedule), nil
}

func (sc *ScheduleUsecase) UpdateSchedule(ctx context.Context, request *model.UpdateScheduleRequest) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	schedule, err := sc.ScheduleRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find schedule: %w", err)
	}
	if schedule == nil {
		return errs.ErrNotFound
	}

	schedule.ShipID = request.ShipID
	schedule.DepartureHarborID = request.DepartureHarborID
	schedule.ArrivalHarborID = request.ArrivalHarborID
	schedule.DepartureDatetime = request.DepartureDatetime
	schedule.ArrivalDatetime = request.ArrivalDatetime
	schedule.Status = request.Status

	if err := sc.ScheduleRepository.Update(tx, schedule); err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sc *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	schedule, err := sc.ScheduleRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}
	if schedule == nil {
		return errs.ErrNotFound
	}

	if err := sc.ScheduleRepository.Delete(tx, schedule); err != nil {
		return fmt.Errorf("failed to delete allocation: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
