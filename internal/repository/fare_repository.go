package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type FareRepository struct {
	Repository[entity.Fare]
}

func NewFareRepository() *FareRepository {
	return &FareRepository{}
}

func (fr *FareRepository) Count(db *gorm.DB) (int64, error) {
	fares := []*entity.Fare{}
	var total int64
	result := db.Find(&fares).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (fr *FareRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Fare, error) {
	fares := []*entity.Fare{}
	result := db.Limit(limit).Offset(offset).Find(&fares)
	if result.Error != nil {
		return nil, result.Error
	}
	return fares, nil
}

func (fr *FareRepository) GetByID(db *gorm.DB, id uint) (*entity.Fare, error) {
	fare := new(entity.Fare)
	result := db.First(&fare, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return fare, result.Error
}

func (fr *FareRepository) GetByManifestAndRoute(db *gorm.DB, manifestID uint, routeID uint) (*entity.Fare, error) {
	fare := new(entity.Fare)
	result := db.Where("manifest_id = ? AND route_id = ?", manifestID, routeID).First(fare)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return fare, result.Error
}
