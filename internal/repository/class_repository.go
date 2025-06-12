package repository

import (
	"errors"
	"eticket-api/internal/entity"

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

func (cr *ClassRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Class, error) {
	classes := []*entity.Class{}
	result := db.Limit(limit).Offset(offset).Find(&classes)
	if result.Error != nil {
		return nil, result.Error
	}
	return classes, nil
}

func (cr *ClassRepository) GetByID(db *gorm.DB, id uint) (*entity.Class, error) {
	class := new(entity.Class)
	result := db.First(&class, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return class, result.Error
}
