package repository

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuotaRepository struct {
	DB *gorm.DB
}

func NewQuotaRepository(db *gorm.DB) *QuotaRepository {
	return &QuotaRepository{DB: db}
}

func (r *QuotaRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Quota{}).Count(&total)
	return total, result.Error
}

func (r *QuotaRepository) Insert(ctx context.Context, conn gotann.Connection, quota *domain.Quota) error {
	result := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Create(quota)
	return result.Error
}

func (r *QuotaRepository) InsertBulk(ctx context.Context, conn gotann.Connection, quotas []*domain.Quota) error {
	result := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Create(&quotas)
	return result.Error
}

func (r *QuotaRepository) Update(ctx context.Context, conn gotann.Connection, quota *domain.Quota) error {
	result := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Save(quota)
	return result.Error
}

func (r *QuotaRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, Quotas []*domain.Quota) error {
	result := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Save(&Quotas)
	return result.Error
}

func (r *QuotaRepository) Delete(ctx context.Context, conn gotann.Connection, quota *domain.Quota) error {
	result := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Select(clause.Associations).Delete(quota)
	return result.Error
}

func (r *QuotaRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Quota, error) {
	Quotas := []*domain.Quota{}
	query := conn.Model(&domain.Quota{}).Preload("Class").
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

func (r *QuotaRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Quota, error) {
	Quota := new(domain.Quota)
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").Clauses(clause.Locking{Strength: "UPDATE"}).First(&Quota, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return Quota, result.Error
}

func (r *QuotaRepository) FindByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) ([]*domain.Quota, error) {
	Quotas := []*domain.Quota{}
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").Clauses(clause.Locking{Strength: "UPDATE"}).Where("schedule_id = ?", scheduleID).Find(&Quotas)
	if result.Error != nil {
		return nil, result.Error // Handle database errors
	}
	return Quotas, nil
}

func (r *QuotaRepository) FindByScheduleIDAndClassID(ctx context.Context, conn gotann.Connection, scheduleID uint, classID uint) (*domain.Quota, error) {
	quota := new(domain.Quota)
	result := conn.Preload("Class").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").Where("schedule_id = ? AND class_id = ?", scheduleID, classID).Clauses(clause.Locking{Strength: "UPDATE"}).First(quota)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return quota, result.Error
}
