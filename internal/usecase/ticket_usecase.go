package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"eticket-api/pkg/utils/tx"
	"fmt"

	"gorm.io/gorm"
)

type TicketUsecase struct {
	Tx                 tx.TxManager
	TicketRepository   *repository.TicketRepository
	ScheduleRepository *repository.ScheduleRepository
	FareRepository     *repository.FareRepository
	SessionRepository  *repository.SessionRepository
}

func NewTicketUsecase(
	tx tx.TxManager,
	ticket_repository *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	fare_repository *repository.FareRepository,
	session_repository *repository.SessionRepository,
) *TicketUsecase {
	return &TicketUsecase{
		Tx:                 tx,
		TicketRepository:   ticket_repository,
		ScheduleRepository: schedule_repository,
		FareRepository:     fare_repository,
		SessionRepository:  session_repository,
	}
}

func (t *TicketUsecase) CreateTicket(ctx context.Context, request *model.WriteTicketRequest) error {
	ticket := mapper.TicketMapper.FromWrite(request)

	if ticket.Status == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return t.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return t.TicketRepository.Create(tx, ticket)
	})
}

func (t *TicketUsecase) GetAllTickets(ctx context.Context, limit, offset int) ([]*model.ReadTicketResponse, int, error) {
	tickets := []*entity.Ticket{}
	var total int64
	err := t.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = t.TicketRepository.Count(tx)
		if err != nil {
			return err
		}
		tickets, err = t.TicketRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	return mapper.TicketMapper.ToModels(tickets), int(total), nil
}

func (t *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	ticket := new(entity.Ticket)

	err := t.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		ticket, err = t.TicketRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}

	if ticket == nil {
		return nil, errors.New("ticket not found")
	}

	return mapper.TicketMapper.ToModel(ticket), nil
}

func (t *TicketUsecase) UpdateTicket(ctx context.Context, id uint, request *model.UpdateTicketRequest) error {
	ticket := mapper.TicketMapper.FromUpdate(request)
	ticket.ID = id

	if ticket.ID == 0 {
		return fmt.Errorf("ticket ID cannot be zero")
	}

	if ticket.PassengerName == nil {
		return fmt.Errorf("passenger name cannot be empty")
	}

	return t.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return t.TicketRepository.Update(tx, ticket)
	})
}

func (t *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {

	return t.Tx.Execute(ctx, func(tx *gorm.DB) error {
		ticket, err := t.TicketRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if ticket == nil {
			return errors.New("ticket not found")
		}
		return t.TicketRepository.Delete(tx, ticket)
	})

}
