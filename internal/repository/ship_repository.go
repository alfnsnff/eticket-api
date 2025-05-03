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

func (shr *ShipRepository) GetAll(db *gorm.DB) ([]*entity.Ship, error) {
	ships := []*entity.Ship{}
	result := db.Find(&ships)
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
