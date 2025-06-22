package repository

import (
	"errors"
	"eticket-api/internal/entity"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FareRepository struct{}

func NewFareRepository() *FareRepository {
	return &FareRepository{}
}

func (ar *FareRepository) Create(db *gorm.DB, fare *entity.Fare) error {
	result := db.Create(fare)
	return result.Error
}

func (ar *FareRepository) Update(db *gorm.DB, fare *entity.Fare) error {
	result := db.Save(fare)
	return result.Error
}

func (ar *FareRepository) Delete(db *gorm.DB, fare *entity.Fare) error {
	result := db.Select(clause.Associations).Delete(fare)
	return result.Error
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

func (fr *FareRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Fare, error) {
	fares := []*entity.Fare{}

	query := db.Preload("Route").
		Preload("Route.DepartureHarbor").
		Preload("Route.ArrivalHarbor").
		Preload("Manifest").
		Preload("Manifest.Class").
		Preload("Manifest.Ship")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("route_id ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&fares).Error
	return fares, err
}

func (fr *FareRepository) GetByID(db *gorm.DB, id uint) (*entity.Fare, error) {
	fare := new(entity.Fare)
	result := db.Preload("Route").
		Preload("Route.DepartureHarbor").
		Preload("Route.ArrivalHarbor").
		Preload("Manifest").
		Preload("Manifest.Class").
		Preload("Manifest.Ship").
		First(&fare, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return fare, result.Error
}

func (fr *FareRepository) GetByManifestAndRoute(db *gorm.DB, manifestID uint, routeID uint) (*entity.Fare, error) {
	fare := new(entity.Fare)
	result := db.Preload("Route").
		Preload("Route.DepartureHarbor").
		Preload("Route.ArrivalHarbor").
		Preload("Manifest").
		Preload("Manifest.Class").
		Preload("Manifest.Ship").
		Where("manifest_id = ? AND route_id = ?", manifestID, routeID).First(fare)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return fare, result.Error
}
