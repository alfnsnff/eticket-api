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

type ClassRepository struct {
	DB *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{DB: db}
}

func (r *ClassRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Class{}).Count(&total)
	return total, result.Error
}

func (r *ClassRepository) Insert(ctx context.Context, conn gotann.Connection, class *domain.Class) error {
	result := conn.Create(class)
	return result.Error
}

func (r *ClassRepository) InsertBulk(ctx context.Context, conn gotann.Connection, classes []*domain.Class) error {
	result := conn.Create(&classes)
	return result.Error
}

func (r *ClassRepository) Update(ctx context.Context, conn gotann.Connection, class *domain.Class) error {
	result := conn.Save(class)
	return result.Error
}

func (r *ClassRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, classes []*domain.Class) error {
	result := conn.Save(&classes)
	return result.Error
}

func (r *ClassRepository) Delete(ctx context.Context, conn gotann.Connection, class *domain.Class) error {
	result := conn.Select(clause.Associations).Delete(class)
	return result.Error
}

func (r *ClassRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Class, error) {
	classes := []*domain.Class{}
	query := conn.Model(&domain.Class{})
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("class_name ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&classes).Error
	return classes, err
}

func (r *ClassRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Class, error) {
	class := new(domain.Class)
	result := conn.First(&class, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return class, result.Error
}
