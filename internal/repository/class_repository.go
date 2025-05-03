package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ClassRepository struct {
	Repository[entity.Class]
}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{}
}

func (cr *ClassRepository) GetAll(db *gorm.DB) ([]*entity.Class, error) {
	classes := []*entity.Class{}
	result := db.Find(&classes)
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
