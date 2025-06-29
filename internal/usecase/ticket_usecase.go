package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type TicketUsecase struct {
	Transactor         *transact.Transactor
	TicketRepository   domain.TicketRepository
	ScheduleRepository domain.ScheduleRepository
	QuotaRepository    domain.QuotaRepository
}

func NewTicketUsecase(

	transactor *transact.Transactor,
	ticket_repository domain.TicketRepository,
	schedule_repository domain.ScheduleRepository,
	quota_reposiotry domain.QuotaRepository,
) *TicketUsecase {
	return &TicketUsecase{

		Transactor:         transactor,
		ScheduleRepository: schedule_repository,
		TicketRepository:   ticket_repository,
		QuotaRepository:    quota_reposiotry,
	}
}
func (uc *TicketUsecase) CreateTicket(ctx context.Context, e *domain.Ticket) error {

	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule, err := uc.ScheduleRepository.FindByID(ctx, tx, e.ScheduleID)
		if err != nil {
			return fmt.Errorf("failed to retrieve schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}

		quota, err := uc.QuotaRepository.FindByScheduleIDAndClassID(ctx, tx, e.ScheduleID, e.ClassID)
		if err != nil {
			return fmt.Errorf("failed to retrieve quota: %w", err)
		}
		if quota == nil {
			return errs.ErrNotFound
		}

		if e.SeatNumber == nil || *e.SeatNumber == "" {
			seat := fmt.Sprintf("%s%d", quota.Class.ClassAlias, quota.Capacity-quota.Quota+1)
			e.SeatNumber = &seat
		}

		ticket := &domain.Ticket{
			TicketCode:      utils.GenerateTicketReferenceID(), // Unique ticket code
			ScheduleID:      e.ScheduleID,
			ClassID:         e.ClassID,
			Type:            quota.Class.Type,
			Price:           quota.Price,
			Address:         e.Address,
			PassengerName:   e.PassengerName,
			PassengerAge:    e.PassengerAge,
			PassengerGender: e.PassengerGender,
			IDType:          e.IDType,
			IDNumber:        e.IDNumber,
			SeatNumber:      e.SeatNumber,
			LicensePlate:    e.LicensePlate,
			IsCheckedIn:     e.IsCheckedIn, // Default value
		}

		if err := uc.TicketRepository.Insert(ctx, tx, ticket); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create ticket: %w", err)
		}

		quota.Quota -= 1

		if err := uc.QuotaRepository.Update(ctx, tx, quota); err != nil {
			return fmt.Errorf("failed to update quota: %w", err)
		}
		return nil
	})
}

func (uc *TicketUsecase) ListTickets(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Ticket, int, error) {
	var err error
	var total int64
	var tickets []*domain.Ticket
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.TicketRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count tickets: %w", err)
		}

		tickets, err = uc.TicketRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all tickets: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets: %w", err)
	}

	return tickets, int(total), nil
}

func (uc *TicketUsecase) ListTicketsByScheduleID(ctx context.Context, schedule_id, limit, offset int, sort, search string) ([]*domain.Ticket, int, error) {
	var err error
	var total int64
	var tickets []*domain.Ticket
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.TicketRepository.CountByScheduleID(ctx, tx, uint(schedule_id))
		if err != nil {
			return fmt.Errorf("failed to count tickets: %w", err)
		}
		tickets, err = uc.TicketRepository.FindByScheduleID(ctx, tx, uint(schedule_id))
		if err != nil {
			return fmt.Errorf("failed to get all tickets: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets by schedule ID: %w", err)
	}

	return tickets, int(total), nil
}

func (uc *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*domain.Ticket, error) {
	var err error
	var ticket *domain.Ticket
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ticket, err = uc.TicketRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get ticket by ID: %w", err)
		}
		if ticket == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}
	return ticket, nil
}

func (uc *TicketUsecase) UpdateTicket(ctx context.Context, e *domain.Ticket) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ticket, err := uc.TicketRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find ticket: %w", err)
		}
		if ticket == nil {
			return errs.ErrNotFound
		}

		ticket.ScheduleID = e.ScheduleID
		ticket.ClassID = e.ClassID
		ticket.Type = e.Type
		ticket.Address = e.Address
		ticket.PassengerName = e.PassengerName
		ticket.PassengerAge = e.PassengerAge
		ticket.PassengerGender = e.PassengerGender
		ticket.IDType = e.IDType
		ticket.IDNumber = e.IDNumber
		ticket.LicensePlate = e.LicensePlate
		ticket.IsCheckedIn = e.IsCheckedIn

		if err := uc.TicketRepository.Update(ctx, tx, ticket); err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		return nil
	})
}

func (uc *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ticket, err := uc.TicketRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get ticket: %w", err)
		}
		if ticket == nil {
			return errs.ErrNotFound
		}

		if err := uc.TicketRepository.Delete(ctx, tx, ticket); err != nil {
			return fmt.Errorf("failed to delete ticket: %w", err)
		}
		return nil
	})
}

func (uc *TicketUsecase) CheckIn(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ticket, err := uc.TicketRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to find ticket: %w", err)
		}
		if ticket == nil {
			return errs.ErrNotFound
		}

		ticket.IsCheckedIn = true

		if err := uc.TicketRepository.Update(ctx, tx, ticket); err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		return nil
	})
}
