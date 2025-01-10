package service

import (
    "fmt"
    "eticket-api/internal/domain"
)

type TicketService struct {
    Repo domain.TicketRepository
}

func (s *TicketService) CreateTicket(ticket *domain.Ticket) error {
    if ticket.Price <= 0 {
        return fmt.Errorf("price must be greater than zero")
    }
    return s.Repo.Create(ticket)
}

func (s *TicketService) GetAllTickets() ([]*domain.Ticket, error) {
    return s.Repo.GetAll()
}

// Additional methods for GetByID, Update, Delete, etc.
