package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type ScheduleUsecase struct {
	Transactor             transact.Transactor
	ClaimSessionRepository domain.ClaimSessionRepository
	ClassRepository        domain.ClassRepository
	ShipRepository         domain.ShipRepository
	ScheduleRepository     domain.ScheduleRepository
	TicketRepository       domain.TicketRepository
}

func NewScheduleUsecase(
	transactor transact.Transactor,
	claim_session_repository domain.ClaimSessionRepository,
	class_repository domain.ClassRepository,
	ship_repository domain.ShipRepository,
	schedule_repository domain.ScheduleRepository,
	ticket_repository domain.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		Transactor:             transactor,
		ClaimSessionRepository: claim_session_repository,
		ClassRepository:        class_repository,
		ShipRepository:         ship_repository,
		ScheduleRepository:     schedule_repository,
		TicketRepository:       ticket_repository,
	}
}

func (uc *ScheduleUsecase) CreateSchedule(ctx context.Context, e *domain.Schedule) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule := &domain.Schedule{
			ShipID:            e.ShipID,
			DepartureHarborID: e.DepartureHarborID,
			ArrivalHarborID:   e.ArrivalHarborID,
			DepartureDatetime: e.DepartureDatetime,
			ArrivalDatetime:   e.ArrivalDatetime,
			Status:            e.Status,
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

func (uc *ScheduleUsecase) ListSchedules(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Schedule, int, error) {
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

	return schedules, int(total), nil
}

func (uc *ScheduleUsecase) ListActiveSchedules(ctx context.Context) ([]*domain.Schedule, error) {
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
	return schedules, nil
}

func (uc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*domain.Schedule, error) {
	var schedule *domain.Schedule
	var err error

	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		// 1. Ambil schedule beserta preloaded quotas
		schedule, err = uc.ScheduleRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}

		// 2. Ambil semua claim session aktif (sudah preload claim items)
		claimSessions, err := uc.ClaimSessionRepository.FindActiveByScheduleID(ctx, tx, schedule.ID)
		if err != nil {
			return fmt.Errorf("failed to get claim sessions: %w", err)
		}

		// 3. Hitung total claimed per class
		claimedByClass := make(map[uint]int)
		for _, session := range claimSessions {
			for _, item := range session.ClaimItems {
				claimedByClass[item.ClassID] += item.Quantity
			}
		}

		// 4. Hitung available quota untuk setiap class
		for _, quota := range schedule.Quotas {
			claimed := claimedByClass[quota.ClassID]
			quota.Quota = max(quota.Quota-claimed, 0)
			fmt.Printf("Schedule ID %d has %d active claim sessions, [Quota:%s:%d]\n", schedule.ID, len(claimSessions), quota.Class.ClassName, quota.Quota)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get schedule with quotas: %w", err)
	}
	return schedule, nil
}

func (uc *ScheduleUsecase) UpdateSchedule(ctx context.Context, e *domain.Schedule) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule, err := uc.ScheduleRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}

		schedule.ShipID = e.ShipID
		schedule.DepartureHarborID = e.DepartureHarborID
		schedule.ArrivalHarborID = e.ArrivalHarborID
		schedule.DepartureDatetime = e.DepartureDatetime
		schedule.ArrivalDatetime = e.ArrivalDatetime
		schedule.Status = e.Status

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
