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
	"time"

	"gorm.io/gorm"
)

type TicketUsecase struct {
	DB                 *gorm.DB
	TicketRepository   *repository.TicketRepository
	ScheduleRepository *repository.ScheduleRepository
	FareRepository     *repository.FareRepository
	SessionRepository  *repository.SessionRepository
}

func NewTicketUsecase(
	db *gorm.DB,
	ticket_repository *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	fare_repository *repository.FareRepository,
	session_repository *repository.SessionRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:                 db,
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

	return tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		return t.TicketRepository.Create(tx, ticket)
	})
}

func (t *TicketUsecase) GetAllTickets(ctx context.Context) ([]*model.ReadTicketResponse, error) {
	tickets := []*entity.Ticket{}

	err := tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		var err error
		tickets, err = t.TicketRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all tickets: %w", err)
	}

	return mapper.TicketMapper.ToModels(tickets), nil
}

func (t *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	ticket := new(entity.Ticket)

	err := tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
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

	return tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		return t.TicketRepository.Update(tx, ticket)
	})
}

func (t *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {

	return tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
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

func (t *TicketUsecase) FillData(ctx context.Context, request *model.FillPassengerDataRequest) (*model.FillPassengerDataResponse, error) {
	if len(request.PassengerData) == 0 {
		return nil, errors.New("invalid request: UserID and passenger data are required")
	}

	_, passengerMap := extractPassengerData(request)

	var updatedIDs []uint
	var failed []model.TicketUpdateFailure

	err := tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		session, err := t.SessionRepository.GetByUUIDWithLock(tx, request.SessionID, true)
		if err != nil {
			return fmt.Errorf("failed to retrieve claim session %s within transaction: %w", request.SessionID, err)
		}
		if session == nil {
			return errors.New("claim session not found")
		}

		now := time.Now()
		if session.ExpiresAt.Before(now) {
			return errors.New("claim session has expired")
		}

		tickets, err := t.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets: %w", err)
		}

		updatedIDs, failed, tickets = t.validateAndUpdateTickets(tickets, passengerMap, now)

		if len(tickets) > 0 {
			err = t.TicketRepository.UpdateBulk(tx, tickets)
			if err != nil {
				return fmt.Errorf("failed to save tickets: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fill passenger data failed: %w", err)
	}

	return &model.FillPassengerDataResponse{
		UpdatedTicketIDs: updatedIDs,
		FailedTickets:    failed,
	}, nil
}

func extractPassengerData(request *model.FillPassengerDataRequest) ([]uint, map[uint]model.PassengerDataInput) {
	ticketIDs := make([]uint, len(request.PassengerData))
	passengerMap := make(map[uint]model.PassengerDataInput)
	for i, data := range request.PassengerData {
		ticketIDs[i] = data.TicketID
		passengerMap[data.TicketID] = data
	}
	return ticketIDs, passengerMap
}

func (t *TicketUsecase) validateAndUpdateTickets(
	tickets []*entity.Ticket,
	dataMap map[uint]model.PassengerDataInput,
	now time.Time,
) ([]uint, []model.TicketUpdateFailure, []*entity.Ticket) {

	retrievedTicketsMap := make(map[uint]*entity.Ticket)
	for _, ticket := range tickets {
		retrievedTicketsMap[ticket.ID] = ticket
	}

	var updatedIDs []uint
	var failed []model.TicketUpdateFailure
	var toUpdate []*entity.Ticket

	for id, data := range dataMap {
		ticket, exists := retrievedTicketsMap[id]
		if !exists {

			failed = append(failed, model.TicketUpdateFailure{TicketID: id, Reason: "Ticket not found in session"})
			continue
		}
		if ticket.Status != "pending_data_entry" {
			failed = append(failed, model.TicketUpdateFailure{TicketID: id, Reason: fmt.Sprintf("Status is %s", ticket.Status)})
			continue
		}
		if data.PassengerName == "" {
			failed = append(failed, model.TicketUpdateFailure{TicketID: id, Reason: "Passenger name required"})
			continue
		}
		if data.IDType == "" {
			failed = append(failed, model.TicketUpdateFailure{TicketID: id, Reason: "ID Type required"})
			continue
		}
		if data.IDNumber == "" {
			failed = append(failed, model.TicketUpdateFailure{TicketID: id, Reason: "ID Number required"})
			continue
		}
		ticket.PassengerName = &data.PassengerName
		ticket.IDType = &data.IDType
		ticket.IDNumber = &data.IDNumber
		ticket.SeatNumber = data.SeatNumber
		ticket.Status = "pending_payment"
		ticket.EntriesAt = &now
		toUpdate = append(toUpdate, ticket)
		updatedIDs = append(updatedIDs, ticket.ID)
	}

	return updatedIDs, failed, toUpdate
}
