package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShipRepository struct{}

func NewShipRepository() *ShipRepository {
	return &ShipRepository{}
}

func (sr *ShipRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Ship{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (sr *ShipRepository) Insert(db *gorm.DB, ship *domain.Ship) error {
	result := db.Create(ship)
	return result.Error
}

func (sr *ShipRepository) InsertBulk(db *gorm.DB, ships []*domain.Ship) error {
	result := db.Create(ships)
	return result.Error
}

func (sr *ShipRepository) Update(db *gorm.DB, ship *domain.Ship) error {
	result := db.Save(ship)
	return result.Error
}

func (sr *ShipRepository) UpdateBulk(db *gorm.DB, ships []*domain.Ship) error {
	result := db.Save(ships)
	return result.Error
}

func (sr *ShipRepository) Delete(db *gorm.DB, ship *domain.Ship) error {
	result := db.Select(clause.Associations).Delete(ship)
	return result.Error
}

func (sr *ShipRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Ship, error) {
	ships := []*domain.Ship{}
	query := db
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("ship_name ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&ships).Error
	return ships, err
}

func (sr *ShipRepository) FindByID(db *gorm.DB, id uint) (*domain.Ship, error) {
	ship := new(domain.Ship)
	result := db.First(&ship, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ship, result.Error
}
