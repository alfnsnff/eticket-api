package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type TicketRepository struct {
	Repository[entity.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (tr *TicketRepository) GetAll(db *gorm.DB) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}
	result := db.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (tr *TicketRepository) GetByID(db *gorm.DB, id uint) (*entity.Ticket, error) {
	ticket := new(entity.Ticket)
	result := db.First(&ticket, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ticket, result.Error
}
func (r *TicketRepository) CountByScheduleClassAndStatuses(db *gorm.DB, scheduleID uint, classID uint, statuses []string) (int64, error) {
	var count int64

	now := time.Now()
	// Define pending statuses that are subject to expiry and linked via Session
	pendingStatuses := []string{"pending_data_entry", "pending_payment"}

	// Use db.Model to specify the table for counting
	query := db.Model(&entity.Ticket{}).
		Where("schedule_id = ? AND class_id = ?", scheduleID, classID).
		// Use LEFT JOIN to include tickets not linked to a session (like 'confirmed')
		// Use the table name 'session' as indicated by your entity and error message
		// Corrected: Use tickets.session_id to match the Ticket entity field name
		Joins("LEFT JOIN claim_session ON ticket.claim_session_id = claim_session.id") // Join tickets FK to session PK

	// --- Build the complex WHERE clause based on the JOINed tables ---
	// We want to count tickets that match the schedule/class AND meet one of these conditions:
	// 1. The status is 'confirmed'
	// OR
	// 2. The status is in the pending statuses
	//    AND the ticket is linked to a Session (session_id IS NOT NULL)
	//    AND the linked Session's expiry time is in the future (session.expires_at > ?)

	query = query.Where(
		// Condition 1: Status is 'confirmed'
		// Use explicit table name 'tickets' for clarity
		db.Where("ticket.status = ?", "confirmed").
			// Condition 2: Pending statuses linked to a non-expired session
			// Use explicit table names for clarity
			// Corrected: Use tickets.session_id to match the Ticket entity field name
			Or(db.Where("ticket.status IN (?) AND ticket.claim_session_id IS NOT NULL AND claim_session.expires_at > ?",
				pendingStatuses, // Parameter for tickets.status IN (?)
				now,             // Parameter for session.expires_at > ?
			)),
	)

	// Execute the count query
	result := query.Count(&count)

	if result.Error != nil {
		// Log this error: Database query failed
		return 0, result.Error
	}

	// Return the count and a nil error if successful
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
	result := db.Where("id IN ?", ids).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (tr *TicketRepository) FindManyBySessionID(db *gorm.DB, sessionID uint) ([]*entity.Ticket, error) {
	tickets := []*entity.Ticket{}
	result := db.Where("claim_session_id = ?", sessionID).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}
