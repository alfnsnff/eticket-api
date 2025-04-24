package usecase

import (
	"errors"
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"fmt"
	"time"
)

type TicketUsecase struct {
	TicketRepository    entities.TicketRepositoryInterface
	ScheduleRepository  entities.ScheduleRepositoryInterface
	ShipClassRepository entities.ShipClassRepositoryInterface
	PriceRepository     entities.PriceRepositoryInterface
}

func NewTicketUsecase(ticketRepository entities.TicketRepositoryInterface,
	scheduleRepository entities.ScheduleRepositoryInterface,
	shipClassRepository entities.ShipClassRepositoryInterface,
	priceRepository entities.PriceRepositoryInterface) TicketUsecase {
	return TicketUsecase{TicketRepository: ticketRepository,
		ScheduleRepository:  scheduleRepository,
		ShipClassRepository: shipClassRepository,
		PriceRepository:     priceRepository}
}

func (s *TicketUsecase) ValidateTicketSelection(req *dto.TicketSelectionRequest) (*dto.TicketSelectionResponse, error) {
	schedule, err := s.ScheduleRepository.GetByID(req.ScheduleID)
	if err != nil || schedule == nil {
		return nil, errors.New("schedule not found")
	}

	var priceIDs []uint
	for _, t := range req.Tickets {
		priceIDs = append(priceIDs, t.PriceID)
	}

	// 1. Get all relevant prices by their IDs
	prices, err := s.PriceRepository.GetByIDs(priceIDs)
	if err != nil {
		return nil, err
	}

	// 2. Build a map to quickly lookup prices by ID
	priceMap := make(map[uint]*entities.Price)
	for _, p := range prices {
		priceMap[p.ID] = p
	}

	var total float32
	var ticketDetails []dto.TicketClassDetailResponse

	for _, t := range req.Tickets {
		price, ok := priceMap[t.PriceID]
		if !ok {
			return nil, fmt.Errorf("invalid price ID: %d", t.PriceID)
		}

		// 3. Get the ShipClass from the price
		shipClass := price.ShipClass
		className := ""
		if shipClass.Class.ID != 0 {
			className = shipClass.Class.Name
		}

		// 4. Check capacity based on schedule + ship_class_id
		booked, err := s.TicketRepository.GetBookedCount(req.ScheduleID, price.ID)
		if err != nil {
			return nil, err
		}

		available := shipClass.Capacity - booked
		if t.Quantity > available {
			return nil, fmt.Errorf("quota exceeded for class: %s", className)
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

	return &dto.TicketSelectionResponse{
		ScheduleID: req.ScheduleID,
		ShipName:   schedule.Ship.Name,
		Datetime:   schedule.Datetime.Format(time.RFC3339),
		Tickets:    ticketDetails,
		Total:      total,
	}, nil
}

// CreateTicket creates a new ticket
func (s *TicketUsecase) CreateTicket(ticket *entities.Ticket) error {
	return s.TicketRepository.Create(ticket)
}

// GetAllTickets retrieves all tickets
func (s *TicketUsecase) GetAllTickets() ([]*entities.Ticket, error) {
	return s.TicketRepository.GetAll()
}

// CreateTicket creates a new ticket
func (s *TicketUsecase) GetBookedCount(scheduleID uint, priceID uint) (int, error) {
	schedule, _ := s.ScheduleRepository.GetByID(scheduleID)
	price, _ := s.PriceRepository.GetByID(priceID)

	if schedule == nil {
		return 0, errors.New("schedule not found")
	}

	if price == nil {
		return 0, errors.New("price not found")
	}

	return s.TicketRepository.GetBookedCount(scheduleID, priceID)
}

// GetTicketByID retrieves a ticket by its ID
func (s *TicketUsecase) GetTicketByID(id uint) (*entities.Ticket, error) {
	ticket, err := s.TicketRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, errors.New("ticket not found")
	}
	return ticket, nil
}

// UpdateTicket updates an existing ticket
func (s *TicketUsecase) UpdateTicket(id uint, ticket *entities.Ticket) error {
	ticket.ID = id

	if ticket.ID == 0 {
		return fmt.Errorf("ship ID cannot be zero")
	}

	if ticket.PassengerName == "" {
		return fmt.Errorf("passenger name cannot be empty")
	}
	return s.TicketRepository.Update(ticket)
}

// DeleteTicket deletes a ticket by its ID
func (s *TicketUsecase) DeleteTicket(id uint) error {
	ticket, err := s.TicketRepository.GetByID(id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}
	return s.TicketRepository.Delete(id)
}
