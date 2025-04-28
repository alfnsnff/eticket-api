package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TicketUsecase struct {
	DB                  *gorm.DB
	TicketRepository    *repository.TicketRepository
	ScheduleRepository  *repository.ScheduleRepository
	ShipClassRepository *repository.ShipClassRepository
	PriceRepository     *repository.PriceRepository
}

func NewTicketUsecase(
	db *gorm.DB,
	ticketRepository *repository.TicketRepository,
	scheduleRepository *repository.ScheduleRepository,
	shipClassRepository *repository.ShipClassRepository,
	priceRepository *repository.PriceRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:                  db,
		TicketRepository:    ticketRepository,
		ScheduleRepository:  scheduleRepository,
		ShipClassRepository: shipClassRepository,
		PriceRepository:     priceRepository,
	}
}

func (s *TicketUsecase) ValidateTicketSelection(ctx context.Context, req *dto.TicketSelectionRequest) (*dto.TicketSelectionResponse, error) {
	var result *dto.TicketSelectionResponse

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(txDB, req.ScheduleID)
		if err != nil || schedule == nil {
			return errors.New("schedule not found")
		}

		var priceIDs []uint
		for _, t := range req.Tickets {
			priceIDs = append(priceIDs, t.PriceID)
		}

		prices, err := s.PriceRepository.GetByIDs(txDB, priceIDs)
		if err != nil {
			return err
		}

		priceMap := make(map[uint]*entities.Price)
		for _, p := range prices {
			priceMap[p.ID] = p
		}

		var total float32
		var ticketDetails []dto.TicketClassDetailResponse

		for _, t := range req.Tickets {
			price, ok := priceMap[t.PriceID]
			if !ok {
				return fmt.Errorf("invalid price ID: %d", t.PriceID)
			}

			shipClass := price.ShipClass
			className := ""
			if shipClass.Class.ID != 0 {
				className = shipClass.Class.Name
			}

			booked, err := s.TicketRepository.GetBookedCount(txDB, req.ScheduleID, price.ID)
			if err != nil {
				return err
			}

			available := shipClass.Capacity - booked
			if t.Quantity > available {
				return fmt.Errorf("quota exceeded for class: %s", className)
			}

			subtotal := price.Price * float32(t.Quantity)
			total += subtotal

			ticketDetails = append(ticketDetails, dto.TicketClassDetailResponse{
				ClassName: className,
				PriceID:   price.ID,
				Price:     price.Price,
				Quantity:  t.Quantity,
				Subtotal:  subtotal,
			})
		}

		result = &dto.TicketSelectionResponse{
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

func (s *TicketUsecase) CreateTicket(ctx context.Context, ticket *entities.Ticket) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.TicketRepository.Create(txDB, ticket)
	})
}

func (s *TicketUsecase) GetAllTickets(ctx context.Context) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		tickets, err = s.TicketRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all tickets: %w", err)
	}

	return tickets, nil
}

func (s *TicketUsecase) GetBookedCount(ctx context.Context, scheduleID uint, priceID uint) (int, error) {
	var count int

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(txDB, scheduleID)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}

		price, err := s.PriceRepository.GetByID(txDB, priceID)
		if err != nil {
			return err
		}
		if price == nil {
			return errors.New("price not found")
		}

		count, err = s.TicketRepository.GetBookedCount(txDB, scheduleID, priceID)
		return err
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*entities.Ticket, error) {
	var ticket *entities.Ticket

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		ticket, err = s.TicketRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}

	if ticket == nil {
		return nil, errors.New("ticket not found")
	}

	return ticket, nil
}

func (s *TicketUsecase) UpdateTicket(ctx context.Context, id uint, ticket *entities.Ticket) error {
	ticket.ID = id

	if ticket.ID == 0 {
		return fmt.Errorf("ticket ID cannot be zero")
	}
	if ticket.PassengerName == "" {
		return fmt.Errorf("passenger name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.TicketRepository.Update(txDB, ticket)
	})
}

func (s *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		ticket, err := s.TicketRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if ticket == nil {
			return errors.New("ticket not found")
		}
		return s.TicketRepository.Delete(txDB, id)
	})
}
