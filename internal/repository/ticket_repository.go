package repository

import (
	"errors"
	"eticket-api/internal/domain"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domain.TicketRepository {
	return &TicketRepository{DB: db}
}

// Create inserts a new ticket into the database
func (r *TicketRepository) Create(ticket *domain.Ticket) error {
	result := r.DB.Create(ticket)
	return result.Error
}

// GetAll retrieves all tickets from the database, including the associated class
func (r *TicketRepository) GetAll() ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	result := r.DB.Preload("Class").Find(&tickets) // Preloads Class relationship
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

// GetByID retrieves a ticket by its ID, including the associated class
func (r *TicketRepository) GetByID(id uint) (*domain.Ticket, error) {
	var ticket domain.Ticket
	result := r.DB.Preload("Class").First(&ticket, id) // Preloads Class and fetches by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ticket is found
	}
	return &ticket, result.Error
}

// Update modifies an existing ticket in the database
func (r *TicketRepository) Update(ticket *domain.Ticket) error {
	// Uses Gorm's Save method to update the ticket
	result := r.DB.Save(ticket)
	return result.Error
}

// Delete removes a ticket from the database by its ID
func (r *TicketRepository) Delete(id uint) error {
	result := r.DB.Delete(&domain.Ticket{}, id) // Deletes the ticket by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no ticket found to delete") // Custom error for non-existent ID
	}
	return nil
}
