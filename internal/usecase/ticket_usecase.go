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

type TicketUsecase struct {
	DB                 *gorm.DB
	TicketRepository   domain.TicketRepository
	ScheduleRepository domain.ScheduleRepository
}

func NewTicketUsecase(
	db *gorm.DB,
	ticket_repository domain.TicketRepository,
	schedule_repository domain.ScheduleRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:                 db,
		ScheduleRepository: schedule_repository,
		TicketRepository:   ticket_repository,
	}
}

func (t *TicketUsecase) CreateTicket(ctx context.Context, request *model.WriteTicketRequest) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule, err := t.ScheduleRepository.FindByID(tx, request.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		return errs.ErrNotFound
	}

	ticket := &domain.Ticket{
		ScheduleID:      request.ScheduleID,
		ClassID:         request.ClassID,
		BookingID:       request.BookingID,
		ClaimSessionID:  request.ClaimSessionID,
		Type:            request.Type,
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

	if err := t.TicketRepository.Insert(tx, ticket); err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) ListTickets(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := t.TicketRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	tickets, err := t.TicketRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	responses := make([]*model.ReadTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = mapper.TicketToResponse(ticket)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (t *TicketUsecase) ListTicketsByScheduleID(ctx context.Context, schedule_id, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := t.TicketRepository.CountByScheduleID(tx, uint(schedule_id))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	tickets, err := t.TicketRepository.FindByScheduleID(tx, uint(schedule_id))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	responses := make([]*model.ReadTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = mapper.TicketToResponse(ticket)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (t *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	ticket, err := t.TicketRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}
	if ticket == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.TicketToResponse(ticket), nil
}

func (t *TicketUsecase) UpdateTicket(ctx context.Context, request *model.UpdateTicketRequest) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Fetch existing allocation
	ticket, err := t.TicketRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find ticket: %w", err)
	}
	if ticket == nil {
		return errs.ErrNotFound
	}

	schedule, err := t.ScheduleRepository.FindByID(tx, request.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		return errs.ErrNotFound
	}

	ticket.ScheduleID = request.ScheduleID
	ticket.ClassID = request.ClassID
	ticket.BookingID = request.BookingID
	ticket.ClaimSessionID = request.ClaimSessionID
	ticket.Type = request.Type
	ticket.Address = request.Address
	ticket.PassengerName = request.PassengerName
	ticket.PassengerAge = request.PassengerAge
	ticket.PassengerGender = request.PassengerGender
	ticket.IDType = request.IDType
	ticket.IDNumber = request.IDNumber
	ticket.SeatNumber = request.SeatNumber
	ticket.LicensePlate = request.LicensePlate
	ticket.IsCheckedIn = request.IsCheckedIn // Default value

	if err := t.TicketRepository.Update(tx, ticket); err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	ticket, err := t.TicketRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return errs.ErrNotFound
	}

	if err := t.TicketRepository.Delete(tx, ticket); err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) CheckIn(ctx context.Context, id uint) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := t.TicketRepository.CheckIn(tx, id); err != nil {
		return fmt.Errorf("failed to check-in ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
