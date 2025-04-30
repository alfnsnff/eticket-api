package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type TicketRepository struct {
	Repository[entities.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

// GetAll retrieves all tickets from the database, including the associated class
func (r *TicketRepository) GetAll(db *gorm.DB) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	result := db.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Booking").
		Preload("Fare.Manifest.Class").
		Preload("Fare.Manifest").
		Preload("Fare").
		Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

// GetByID retrieves a ticket by its ID, including the associated class
func (r *TicketRepository) GetByID(db *gorm.DB, id uint) (*entities.Ticket, error) {
	var ticket entities.Ticket
	result := db.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Fare.Manifest.Class").
		Preload("Fare.Manifest").
		Preload("Fare").
		Preload("Booking").
		First(&ticket, id) // Preloads Class and fetches by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ticket is found
	}
	return &ticket, result.Error
}

func (r *TicketRepository) GetBookedCount(db *gorm.DB, scheduleID uint, priceID uint) (int, error) {
	var count int64
	err := db.Table("ticket").
		Where("schedule_id = ? AND price_id = ?", scheduleID, priceID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
