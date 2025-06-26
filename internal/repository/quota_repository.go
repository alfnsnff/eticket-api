package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuotaRepository struct{}

func NewQuotaRepository() *QuotaRepository {
	return &QuotaRepository{}
}

func (ar *QuotaRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Quota{}).Count(&total)
	return total, result.Error
}

func (ar *QuotaRepository) Insert(db *gorm.DB, Quota *domain.Quota) error {
	result := db.Create(Quota)
	return result.Error
}

func (ar *QuotaRepository) InsertBulk(db *gorm.DB, Quotas []*domain.Quota) error {
	result := db.Create(&Quotas)
	return result.Error
}

func (ar *QuotaRepository) Update(db *gorm.DB, Quota *domain.Quota) error {
	result := db.Save(Quota)
	return result.Error
}

func (ar *QuotaRepository) UpdateBulk(db *gorm.DB, Quotas []*domain.Quota) error {
	result := db.Save(&Quotas)
	return result.Error
}

func (ar *QuotaRepository) Delete(db *gorm.DB, Quota *domain.Quota) error {
	result := db.Select(clause.Associations).Delete(Quota)
	return result.Error
}

func (ar *QuotaRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Quota, error) {
	Quotas := []*domain.Quota{}
	query := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship")
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("schedule_id ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&Quotas).Error
	return Quotas, err
}

func (ar *QuotaRepository) FindByID(db *gorm.DB, id uint) (*domain.Quota, error) {
	Quota := new(domain.Quota)
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").First(&Quota, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return Quota, result.Error
}

func (ar *QuotaRepository) FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*domain.Quota, error) {
	Quotas := []*domain.Quota{}
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").Where("schedule_id = ?", scheduleID).Find(&Quotas)
	if result.Error != nil {
		return nil, result.Error // Handle database errors
	}
	return Quotas, nil
}

func (ar *QuotaRepository) FindByScheduleIDAndClassID(db *gorm.DB, scheduleID uint, classID uint) (*domain.Quota, error) {
	Quota := new(domain.Quota)
	result := db.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").Where("schedule_id = ? AND class_id = ?", scheduleID, classID).Clauses(clause.Locking{Strength: "UPDATE"}).First(Quota)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return Quota, result.Error
}
