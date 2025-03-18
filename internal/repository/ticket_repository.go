package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) entities.TicketRepositoryInterface {
	return &TicketRepository{DB: db}
}

// Create inserts a new ticket into the database
func (r *TicketRepository) Create(ticket *entities.Ticket) error {
	result := r.DB.Create(ticket)
	return result.Error
}

// GetAll retrieves all tickets from the database, including the associated class
func (r *TicketRepository) GetAll() ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	result := r.DB.
		Preload("Booking.Schedule.Route.DepartureHarbor").
		Preload("Booking.Schedule.Route.ArrivalHarbor").
		Preload("Booking.Schedule.Ship").
		Preload("Class.Route.DepartureHarbor"). // Preload from Class
		Preload("Class.Route.ArrivalHarbor").   // Preload from Class
		Preload("Class").
		Preload("Booking").
		Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

// GetByID retrieves a ticket by its ID, including the associated class
func (r *TicketRepository) GetByID(id uint) (*entities.Ticket, error) {
	var ticket entities.Ticket
	result := r.DB.Preload("Booking.Schedule.Route.DepartureHarbor").
		Preload("Booking.Schedule.Route.ArrivalHarbor").
		Preload("Booking.Schedule.Ship").
		Preload("Class.Route.DepartureHarbor"). // Preload from Class
		Preload("Class.Route.ArrivalHarbor").   // Preload from Class
		Preload("Class").
		Preload("Booking").First(&ticket, id) // Preloads Class and fetches by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ticket is found
	}
	return &ticket, result.Error
}

// Update modifies an existing ticket in the database
func (r *TicketRepository) Update(ticket *entities.Ticket) error {
	// Uses Gorm's Save method to update the ticket
	result := r.DB.Save(ticket)
	return result.Error
}

// Delete removes a ticket from the database by its ID
func (r *TicketRepository) Delete(id uint) error {
	result := r.DB.Delete(&entities.Ticket{}, id) // Deletes the ticket by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no ticket found to delete") // Custom error for non-existent ID
	}
	return nil
}
