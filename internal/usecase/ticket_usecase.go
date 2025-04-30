package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
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
}

func NewTicketUsecase(
	db *gorm.DB,
	ticket_repository *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	fare_repository *repository.FareRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:                 db,
		TicketRepository:   ticket_repository,
		ScheduleRepository: schedule_repository,
		FareRepository:     fare_repository,
	}
}

func (s *TicketUsecase) ValidateTicketSelection(ctx context.Context, req *model.TicketSelectionRequest) (*model.TicketSelectionResponse, error) {
	var result *model.TicketSelectionResponse

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(tx, req.ScheduleID)
		if err != nil || schedule == nil {
			return errors.New("schedule not found")
		}

		var fareIDs []uint
		for _, t := range req.Tickets {
			fareIDs = append(fareIDs, t.FareID)
		}

		fares, err := s.FareRepository.GetByIDs(tx, fareIDs)
		if err != nil {
			return err
		}

		fareMap := make(map[uint]*entities.Fare)
		for _, f := range fares {
			fareMap[f.ID] = f
		}

		var total float32
		var ticketDetails []model.TicketClassDetailResponse

		for _, t := range req.Tickets {
			fare, ok := fareMap[t.FareID]
			if !ok {
				return fmt.Errorf("invalid price ID: %d", t.FareID)
			}

			manifest := fare.Manifest
			className := ""
			if manifest.Class.ID != 0 {
				className = manifest.Class.Name
			}

			booked, err := s.TicketRepository.GetBookedCount(tx, req.ScheduleID, fare.ID)
			if err != nil {
				return err
			}

			available := manifest.Capacity - booked
			if t.Quantity > available {
				return fmt.Errorf("quota exceeded for class: %s", className)
			}

			subtotal := fare.Price * float32(t.Quantity)
			total += subtotal

			ticketDetails = append(ticketDetails, model.TicketClassDetailResponse{
				ClassName: className,
				FareID:    fare.ID,
				Price:     fare.Price,
				Quantity:  t.Quantity,
				Subtotal:  subtotal,
			})
		}

		result = &model.TicketSelectionResponse{
			ScheduleID: req.ScheduleID,
			ShipName:   schedule.Ship.Name,
			Datetime:   schedule.Datetime.Format(time.RFC3339),
			Tickets:    ticketDetails,
			Total:      total,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TicketUsecase) CreateTicket(ctx context.Context, request *model.WriteTicketRequest) error {
	ticket := mapper.ToTicketEntity(request)

	if ticket.BookingID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}

	if ticket.PassengerName == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.TicketRepository.Create(tx, ticket)
	})
}

func (s *TicketUsecase) GetAllTickets(ctx context.Context) ([]*model.ReadTicketResponse, error) {
	tickets := []*entities.Ticket{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		tickets, err = s.TicketRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all tickets: %w", err)
	}

	return mapper.ToTicketsModel(tickets), nil
}

func (s *TicketUsecase) GetBookedCount(ctx context.Context, request *model.CountBookedTicketRequest) (int, error) {
	var count int

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(tx, request.ScheduleID)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}

		price, err := s.FareRepository.GetByID(tx, request.FareID)
		if err != nil {
			return err
		}
		if price == nil {
			return errors.New("price not found")
		}

		count, err = s.TicketRepository.GetBookedCount(tx, request.ScheduleID, request.FareID)
		return err
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	ticket := new(entities.Ticket)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		ticket, err = s.TicketRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}

	if ticket == nil {
		return nil, errors.New("ticket not found")
	}

	return mapper.ToTicketModel(ticket), nil
}

func (s *TicketUsecase) UpdateTicket(ctx context.Context, id uint, request *model.WriteTicketRequest) error {

	ticket := mapper.ToTicketEntity(request)
	ticket.ID = id

	if ticket.ID == 0 {
		return fmt.Errorf("ticket ID cannot be zero")
	}
	if ticket.PassengerName == "" {
		return fmt.Errorf("passenger name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.TicketRepository.Update(tx, ticket)
	})
}

func (s *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		ticket, err := s.TicketRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if ticket == nil {
			return errors.New("ticket not found")
		}
		return s.TicketRepository.Delete(tx, ticket)
	})
}
