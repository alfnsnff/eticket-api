package repository

import (
    "eticket-api/internal/domain"
    "gorm.io/gorm"
)

type TicketRepositoryImpl struct {
    DB *gorm.DB
}

// Create inserts a new ticket into the database
func (r *TicketRepositoryImpl) Create(ticket *domain.Ticket) error {
    result := r.DB.Create(ticket)
    return result.Error
}

// GetAll retrieves all tickets from the database
func (r *TicketRepositoryImpl) GetAll() ([]*domain.Ticket, error) {
    var tickets []*domain.Ticket
    result := r.DB.Find(&tickets)
    if result.Error != nil {
        return nil, result.Error
    }
    return tickets, nil
}

