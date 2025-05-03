package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

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

func (tr *TicketRepository) CountByScheduleClassAndStatuses(db *gorm.DB, scheduleID uint, classID uint, statuses []string) (int64, error) {
	ticket := new(entity.Ticket)
	result := db.Find(&ticket, "schedule_id = ? AND class_id = ? AND status IN ?", scheduleID, classID, statuses)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
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
