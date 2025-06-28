package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
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
func (uc *TicketUsecase) CreateTicket(ctx context.Context, request *model.WriteTicketRequest) error {

	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		schedule, err := uc.ScheduleRepository.FindByID(ctx, tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("failed to retrieve schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}

		quota, err := uc.QuotaRepository.FindByScheduleIDAndClassID(ctx, tx, request.ScheduleID, request.ClassID)
		if err != nil {
			return fmt.Errorf("failed to retrieve quota: %w", err)
		}
		if quota == nil {
			return errs.ErrNotFound
		}

		if request.SeatNumber == nil || *request.SeatNumber == "" {
			seat := fmt.Sprintf("%s%d", quota.Class.ClassAlias, quota.Capacity-quota.Quota+1)
			request.SeatNumber = &seat
		}

		ticket := &domain.Ticket{
			TicketCode:      utils.GenerateTicketReferenceID(), // Unique ticket code
			ScheduleID:      request.ScheduleID,
			ClassID:         request.ClassID,
			Type:            quota.Class.Type,
			Price:           quota.Price,
			Address:         request.Address,
			PassengerName:   request.PassengerName,
			PassengerAge:    request.PassengerAge,
			PassengerGender: request.PassengerGender,
			IDType:          request.IDType,
			IDNumber:        request.IDNumber,
			SeatNumber:      request.SeatNumber,
			LicensePlate:    request.LicensePlate,
			IsCheckedIn:     request.IsCheckedIn, // Default value
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

func (uc *TicketUsecase) ListTickets(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
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

	responses := make([]*model.ReadTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = mapper.TicketToResponse(ticket)
	}

	return responses, int(total), nil
}

func (uc *TicketUsecase) ListTicketsByScheduleID(ctx context.Context, schedule_id, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
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
	responses := make([]*model.ReadTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = mapper.TicketToResponse(ticket)
	}

	return responses, int(total), nil
}

func (uc *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	var err error
	var ticket *domain.Ticket
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ticket, err := uc.TicketRepository.FindByID(ctx, tx, id)
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
	return mapper.TicketToResponse(ticket), nil
}

func (uc *TicketUsecase) UpdateTicket(ctx context.Context, request *model.UpdateTicketRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ticket, err := uc.TicketRepository.FindByID(ctx, tx, request.ID)
		if err != nil {
			return fmt.Errorf("failed to find ticket: %w", err)
		}
		if ticket == nil {
			return errs.ErrNotFound
		}

		ticket.ScheduleID = request.ScheduleID
		ticket.ClassID = request.ClassID
		ticket.Type = request.Type
		ticket.Address = request.Address
		ticket.PassengerName = request.PassengerName
		ticket.PassengerAge = request.PassengerAge
		ticket.PassengerGender = request.PassengerGender
		ticket.IDType = request.IDType
		ticket.IDNumber = request.IDNumber
		ticket.LicensePlate = request.LicensePlate
		ticket.IsCheckedIn = request.IsCheckedIn

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
