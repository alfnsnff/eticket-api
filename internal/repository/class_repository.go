package repository

import (
	"errors"
	"eticket-api/internal/entity"
	"strings"

	"gorm.io/gorm"
)

type ClassRepository struct {
	Repository[entity.Class]
	DB *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{DB: db}
}

func (cr *ClassRepository) Count(db *gorm.DB) (int64, error) {
	classes := []*entity.Class{}
	var total int64
	result := db.Find(&classes).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (cr *ClassRepository) Createtest(entity *entity.Class) error {
	return cr.DB.Create(entity).Error
}

func (cr *ClassRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Class, error) {
	classes := []*entity.Class{}

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

func (cr *ClassRepository) GetByID(db *gorm.DB, id uint) (*entity.Class, error) {
	class := new(entity.Class)
	result := db.First(&class, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return class, result.Error
}
