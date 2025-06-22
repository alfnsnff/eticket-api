package ticket

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type TicketUsecase struct {
	DB               *gorm.DB
	TicketRepository TicketRepository
}

func NewTicketUsecase(
	db *gorm.DB,
	ticket_repository TicketRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:               db,
		TicketRepository: ticket_repository,
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

	ticket := &entity.Ticket{
		ScheduleID:      request.ScheduleID,
		ClassID:         request.ClassID,
		BookingID:       request.BookingID,
		ClaimSessionID:  request.ClaimSessionID,
		Type:            request.Type,
		Price:           request.Price,
		PassengerName:   request.PassengerName,
		PassengerAge:    request.PassengerAge,
		PassengerGender: request.PassengerGender,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		SeatNumber:      request.SeatNumber,
		LicensePlate:    request.LicensePlate,
	}

	if err := t.TicketRepository.Create(tx, ticket); err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) GetAllTickets(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := t.TicketRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	tickets, err := t.TicketRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadTicketResponses(tickets), int(total), nil
}

func (t *TicketUsecase) GetAllTicketsByScheduleID(ctx context.Context, schedule_id, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	var total int64
	countt := tx.Model(&entity.Ticket{}).Where("schedule_id = ?", schedule_id)
	if search != "" {
		search = "%" + search + "%"
		countt = countt.Where("passenger_name ILIKE ?", search)
	}
	if err := countt.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	tickets, err := t.TicketRepository.GetByScheduleID(tx, schedule_id, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadTicketResponses(tickets), int(total), nil
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

	ticket, err := t.TicketRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}
	if ticket == nil {
		return nil, errors.New("ticket not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadTicketResponse(ticket), nil
}

func (t *TicketUsecase) UpdateTicket(ctx context.Context, request *model.UpdateTicketRequest) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	ticket, err := t.TicketRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}

	ticket.ScheduleID = request.ScheduleID
	ticket.ClassID = request.ClassID
	ticket.BookingID = &request.BookingID
	ticket.ClaimSessionID = &request.ClaimSessionID
	ticket.Type = request.Type
	ticket.Price = request.Price
	ticket.PassengerName = &request.PassengerName
	ticket.PassengerAge = &request.PassengerAge
	ticket.PassengerGender = &request.PassengerGender
	ticket.IDType = &request.IDType
	ticket.IDNumber = &request.IDNumber
	ticket.SeatNumber = request.SeatNumber
	ticket.LicensePlate = request.LicensePlate

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
		} else {
			tx.Rollback()
		}
	}()

	ticket, err := t.TicketRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
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
		} else {
			tx.Rollback()
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
