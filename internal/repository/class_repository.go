package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClassRepository struct{}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{}
}

func (cr *ClassRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Class{}).Count(&total)
	return total, result.Error
}

func (ar *ClassRepository) Insert(db *gorm.DB, class *domain.Class) error {
	result := db.Create(class)
	return result.Error
}

func (cr *ClassRepository) InsertBulk(db *gorm.DB, classes []*domain.Class) error {
	result := db.Create(&classes)
	return result.Error
}

func (cr *ClassRepository) Update(db *gorm.DB, class *domain.Class) error {
	result := db.Save(class)
	return result.Error
}

func (cr *ClassRepository) UpdateBulk(db *gorm.DB, classes []*domain.Class) error {
	result := db.Save(&classes)
	return result.Error
}

func (cr *ClassRepository) Delete(db *gorm.DB, class *domain.Class) error {
	result := db.Select(clause.Associations).Delete(class)
	return result.Error
}

func (cr *ClassRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Class, error) {
	classes := []*domain.Class{}
	query := db
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

func (cr *ClassRepository) FindByID(db *gorm.DB, id uint) (*domain.Class, error) {
	class := new(domain.Class)
	result := db.First(&class, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return class, result.Error
}
