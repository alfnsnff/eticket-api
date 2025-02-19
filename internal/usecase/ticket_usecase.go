package usecase

import (
	"errors"
	"eticket-api/internal/domain"
)

type TicketUsecase struct {
	TicketRepository domain.TicketRepository
}

func NewTicketUsecase(ticketRepository domain.TicketRepository) TicketUsecase {
	return TicketUsecase{TicketRepository: ticketRepository}
}

// CreateTicket creates a new ticket
func (s *TicketUsecase) CreateTicket(ticket *domain.Ticket) error {
	return s.TicketRepository.Create(ticket)
}

// GetAllTickets retrieves all tickets
func (s *TicketUsecase) GetAllTickets() ([]*domain.Ticket, error) {
	return s.TicketRepository.GetAll()
}

// GetTicketByID retrieves a ticket by its ID
func (s *TicketUsecase) GetTicketByID(id uint) (*domain.Ticket, error) {
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
func (s *TicketUsecase) UpdateTicket(ticket *domain.Ticket) error {
	if ticket.ID == 0 {
		return errors.New("ticket ID is required for update")
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
