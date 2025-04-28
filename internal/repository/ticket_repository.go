package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
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
	result := r.DB.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Booking").
		Preload("Price.ShipClass.Class").
		Preload("Price.ShipClass").
		Preload("Price").
		Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

// GetByID retrieves a ticket by its ID, including the associated class
func (r *TicketRepository) GetByID(id uint) (*entities.Ticket, error) {
	var ticket entities.Ticket
	result := r.DB.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Price.ShipClass.Class").
		Preload("Price.ShipClass").
		Preload("Price").
		Preload("Booking").
		First(&ticket, id) // Preloads Class and fetches by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ticket is found
	}
	return &ticket, result.Error
}

func (r *TicketRepository) GetBookedCount(scheduleID uint, priceID uint) (int, error) {
	var count int64
	err := r.DB.Table("ticket").
		Where("schedule_id = ? AND price_id = ?", scheduleID, priceID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
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
