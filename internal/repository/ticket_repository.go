package repository

import (
	"errors"
	"eticket-api/internal/entity"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TicketRepository struct {
	Repository[entity.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (tr *TicketRepository) Count(db *gorm.DB) (int64, error) {
	tickets := []*entity.Ticket{}
	var total int64
	result := db.Find(&tickets).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (tr *TicketRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}

	query := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("passenger_name ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&tickets).Error
	return tickets, err
}

func (tr *TicketRepository) GetBySchedulseID(db *gorm.DB, id, limit, offset int, sort, search string) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}

	query := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("passenger_name ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Where("schedule_id = ?", id).Limit(limit).Offset(offset).Find(&tickets).Error
	return tickets, err
}

func (tr *TicketRepository) GetByID(db *gorm.DB, id uint) (*entity.Ticket, error) {
	ticket := new(entity.Ticket)
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").First(&ticket, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ticket, result.Error
}

func (tr *TicketRepository) GetByBookingID(db *gorm.DB, id uint) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Where("booking_id = ?", id).Find(&tickets)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return tickets, result.Error
}

// func (tr *TicketRepository) GetBySchedulseID(db *gorm.DB, id uint) ([]*entity.Ticket, error) {
// 	tickets := []*entity.Ticket{}
// 	result := db.Preload("Class").
// 		Preload("Schedule").
// 		Preload("Schedule.Ship").
// 		Preload("Schedule.Route").
// 		Preload("Schedule.Route.DepartureHarbor").
// 		Preload("Schedule.Route.ArrivalHarbor").
// 		Where("schedule_id = ?", id).Find(&tickets)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return nil, nil
// 	}
// 	return tickets, result.Error
// }

func (r *TicketRepository) CountByScheduleClassAndStatuses(db *gorm.DB, scheduleID uint, classID uint, statuses []string) (int64, error) {
	var count int64
	now := time.Now()
	pendingStatuses := []string{"pending_data_entry", "pending_payment"}

	query := db.Model(&entity.Ticket{}).
		Joins("LEFT JOIN claim_session ON ticket.claim_session_id = claim_session.id").
		Where("ticket.schedule_id = ? AND ticket.class_id = ?", scheduleID, classID).
		Where(
			db.Where("ticket.status = ?", "confirmed").
				Or("ticket.status IN ? AND ticket.claim_session_id IS NOT NULL AND claim_session.expires_at > ?", pendingStatuses, now),
		)

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (tr *TicketRepository) CreateBulk(db *gorm.DB, tickets []*entity.Ticket) error {
	result := db.Create(&tickets)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (tr *TicketRepository) UpdateBulk(db *gorm.DB, tickets []*entity.Ticket) error {
	result := db.Save(&tickets)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (tr *TicketRepository) FindManyByIDs(db *gorm.DB, ids []uint) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Where("id IN ?", ids).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (tr *TicketRepository) FindManyBySessionID(db *gorm.DB, sessionID uint) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Where("claim_session_id = ?", sessionID).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (r *TicketRepository) CancelManyBySessionID(db *gorm.DB, sessionID uint) error {
	result := db.Model(&entity.Ticket{}).
		Where("claim_session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":           "cancelled",
			"claim_session_id": nil,
		})

	// Check for database errors during the update operation
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *TicketRepository) Paid(db *gorm.DB, id uint) error {
	return db.Model(&entity.Ticket{}).
		Where("id = ?", id).
		Update("status", "paid").
		Error
}

func (r *TicketRepository) CheckIn(db *gorm.DB, id uint) error {
	return db.Model(&entity.Ticket{}).
		Where("id = ?", id).
		Update("status", "checkin").
		Error
}
