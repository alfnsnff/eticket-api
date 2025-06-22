package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClassRepository struct{}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{}
}

func (ar *ClassRepository) Create(db *gorm.DB, class *domain.Class) error {
	result := db.Create(class)
	return result.Error
}

func (ar *ClassRepository) Update(db *gorm.DB, class *domain.Class) error {
	result := db.Save(class)
	return result.Error
}

func (ar *ClassRepository) Delete(db *gorm.DB, class *domain.Class) error {
	result := db.Select(clause.Associations).Delete(class)
	return result.Error
}

func (cr *ClassRepository) Count(db *gorm.DB) (int64, error) {
	classes := []*domain.Class{}
	var total int64
	result := db.Find(&classes).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (cr *ClassRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Class, error) {
	classes := []*domain.Class{}

	query := db

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("class_name ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&classes).Error
	return classes, err
}

func (cr *ClassRepository) GetByID(db *gorm.DB, id uint) (*domain.Class, error) {
	class := new(domain.Class)
	result := db.First(&class, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return class, result.Error
}
