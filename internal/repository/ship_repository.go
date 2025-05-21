package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ShipRepository struct {
	Repository[entity.Ship]
}

func NewShipRepository() *ShipRepository {
	return &ShipRepository{}
}

func (shr *ShipRepository) Count(db *gorm.DB) (int64, error) {
	ships := []*entity.Ship{}
	var total int64
	result := db.Find(&ships).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (shr *ShipRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Ship, error) {
	ships := []*entity.Ship{}
	result := db.Limit(limit).Offset(offset).Find(&ships)
	if result.Error != nil {
		return nil, result.Error
	}
	return ships, nil
}

func (shr *ShipRepository) GetByID(db *gorm.DB, id uint) (*entity.Ship, error) {
	ship := new(entity.Ship)
	result := db.First(&ship, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ship, result.Error
}
