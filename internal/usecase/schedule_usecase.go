package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"eticket-api/pkg/gotann"
	"fmt"
)

type ScheduleUsecase struct {
	Transactor         *transact.Transactor
	ClassRepository    domain.ClassRepository
	ShipRepository     domain.ShipRepository
	ScheduleRepository domain.ScheduleRepository
	TicketRepository   domain.TicketRepository
}

func NewScheduleUsecase(

	transactor *transact.Transactor,
	class_repository domain.ClassRepository,
	ship_repository domain.ShipRepository,
	schedule_repository domain.ScheduleRepository,
	ticket_repository domain.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{

		Transactor:         transactor,
		ClassRepository:    class_repository,
		ShipRepository:     ship_repository,
		ScheduleRepository: schedule_repository,
		TicketRepository:   ticket_repository,
	}
}

func (uc *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.WriteScheduleRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule := &domain.Schedule{
			ShipID:            request.ShipID,
			DepartureHarborID: request.DepartureHarborID,
			ArrivalHarborID:   request.ArrivalHarborID,
			DepartureDatetime: request.DepartureDatetime,
			ArrivalDatetime:   request.ArrivalDatetime,
			Status:            request.Status,
		}
		if err := uc.ScheduleRepository.Insert(ctx, tx, schedule); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create schedule: %w", err)
		}
		return nil
	})
}

func (uc *ScheduleUsecase) ListSchedules(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadScheduleResponse, int, error) {
	var err error
	var total int64
	schedules := []*domain.Schedule{}
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.ScheduleRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count schedules: %w", err)
		}

		schedules, err = uc.ScheduleRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all schedules: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list schedules: %w", err)
	}
	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = mapper.ScheduleToResponse(schedule)
	}
	return responses, int(total), nil
}

func (uc *ScheduleUsecase) ListActiveSchedules(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	var err error
	var schedules []*domain.Schedule
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedules, err = uc.ScheduleRepository.FindActiveSchedules(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get all active schedules: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to list active schedules: %w", err)
	}
	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = mapper.ScheduleToResponse(schedule)
	}
	return responses, nil
}

func (uc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	var err error
	var schedule *domain.Schedule
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule, err = uc.ScheduleRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	return mapper.ScheduleToResponse(schedule), nil
}

func (uc *ScheduleUsecase) UpdateSchedule(ctx context.Context, request *model.UpdateScheduleRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule, err := uc.ScheduleRepository.FindByID(ctx, tx, request.ID)
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

		if err := uc.ScheduleRepository.Update(ctx, tx, schedule); err != nil {
			return fmt.Errorf("failed to update schedule: %w", err)
		}
		return nil
	})
}

func (uc *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule, err := uc.ScheduleRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}

		if err := uc.ScheduleRepository.Delete(ctx, tx, schedule); err != nil {
			return fmt.Errorf("failed to delete allocation: %w", err)
		}
		return nil
	})
}
