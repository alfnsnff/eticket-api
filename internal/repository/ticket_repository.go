package repository

import (
	"context"
	"errors"
	enum "eticket-api/internal/common/enums"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

func (r *TicketRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Ticket{}).Count(&total)
	return total, result.Error
}

func (r *TicketRepository) CountByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) (int64, error) {
	var total int64
	result := conn.Model(&domain.Ticket{}).Where("schedule_id = ?", scheduleID).Count(&total)
	return total, result.Error
}

func (r *TicketRepository) CountByScheduleIDAndClassIDWithStatus(ctx context.Context, conn gotann.Connection, scheduleID uint, classID uint) (int64, error) {
	var count int64
	now := time.Now()

	query := conn.Model(&domain.Ticket{}).
		Joins("LEFT JOIN claim_session ON ticket.claim_session_id = claim_session.id").
		Where("ticket.schedule_id = ? AND ticket.class_id = ?", scheduleID, classID).
		Where("ticket.claim_session_id IS NOT NULL").
		Where(`
            (claim_session.status IN ?) OR 
            (claim_session.expires_at > ? AND claim_session.status IN ?)
        `, enum.GetSuccessClaimSessionStatuses(), now, enum.GetPendingClaimSessionStatuses())

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (r *TicketRepository) Insert(ctx context.Context, conn gotann.Connection, ticket *domain.Ticket) error {
	result := conn.Create(ticket)
	return result.Error
}

func (r *TicketRepository) InsertBulk(ctx context.Context, conn gotann.Connection, tickets []*domain.Ticket) error {
	result := conn.Create(&tickets)
	return result.Error
}

func (r *TicketRepository) Update(ctx context.Context, conn gotann.Connection, ticket *domain.Ticket) error {
	result := conn.Save(ticket)
	return result.Error
}

func (r *TicketRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, tickets []*domain.Ticket) error {
	result := conn.Save(&tickets)
	return result.Error
}

func (r *TicketRepository) Delete(ctx context.Context, conn gotann.Connection, ticket *domain.Ticket) error {
	result := conn.Select(clause.Associations).Delete(ticket)
	return result.Error
}

func (r *TicketRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Ticket, error) {
	tickets := []*domain.Ticket{}
	query := conn.Model(&domain.Ticket{}).Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor")
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("passenger_name ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Ticket, error) {
	ticket := new(domain.Ticket)
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Booking").
		First(&ticket, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ticket, result.Error
}

func (r *TicketRepository) FindByIDs(ctx context.Context, conn gotann.Connection, ids []uint) ([]*domain.Ticket, error) {
	tickets := []*domain.Ticket{}
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Booking").
		Where("id IN ?", ids).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (r *TicketRepository) FindByBookingID(ctx context.Context, conn gotann.Connection, bookingID uint) ([]*domain.Ticket, error) {
	tickets := []*domain.Ticket{}
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Booking").
		Where("booking_id = ?", bookingID).
		Find(&tickets)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return tickets, result.Error
}

func (r *TicketRepository) FindByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) ([]*domain.Ticket, error) {
	tickets := []*domain.Ticket{}
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Booking").
		Where("schedule_id = ?", scheduleID).
		Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (r *TicketRepository) FindByClaimSessionID(ctx context.Context, conn gotann.Connection, sessionID uint) ([]*domain.Ticket, error) {
	tickets := []*domain.Ticket{}
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Booking").
		Where("claim_session_id = ?", sessionID).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (r *TicketRepository) CheckIn(ctx context.Context, conn gotann.Connection, id uint) error {
	return conn.Model(&domain.Ticket{}).
		Where("id = ?", id).
		Update("status", "checkin").
		Error
}
