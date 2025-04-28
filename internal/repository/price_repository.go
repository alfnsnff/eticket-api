package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type PriceRepository struct {
	DB *gorm.DB
}

func NewPriceRepository() *PriceRepository {
	return &PriceRepository{}
}

// Create inserts a new route into the database
func (p *PriceRepository) Create(db *gorm.DB, price *entities.Price) error {
	result := db.Create(price)
	return result.Error
}

// GetAll retrieves all routes from the database
func (p *PriceRepository) GetAll(db *gorm.DB) ([]*entities.Price, error) {
	var Prices []*entities.Price
	result := db.Find(&Prices) // Corrected Preload
	if result.Error != nil {
		return nil, result.Error
	}
	return Prices, nil
}

// GetByID retrieves a route by its ID
func (p *PriceRepository) GetByID(db *gorm.DB, id uint) (*entities.Price, error) {
	var Price entities.Price
	result := db.First(&Price, id) // Fetches the route by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return &Price, result.Error
}

// GetByID retrieves a route by its ID
func (p *PriceRepository) GetByIDs(db *gorm.DB, ids []uint) ([]*entities.Price, error) {
	var Prices []*entities.Price

	result := db.Where("id IN ?", ids).Preload("ShipClass").
		Preload("ShipClass.Class").Find(&Prices) // Fetches the route by ID

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return Prices, result.Error
}

// PriceRepository.go
func (r *PriceRepository) GetByRouteID(db *gorm.DB, routeID uint) ([]*entities.Price, error) {
	var Prices []*entities.Price
	result := db.
		Preload("ShipClass").
		Preload("ShipClass.Class"). // for ClassName
		Where("route_id = ?", routeID).
		Find(&Prices)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return Prices, result.Error
}

// Update modifies an existing route in the database
func (p *PriceRepository) Update(db *gorm.DB, route *entities.Price) error {
	// Uses Gorm's Save method to update the Route
	result := db.Save(route)
	return result.Error
}

// Delete removes a route from the database by its ID
func (p *PriceRepository) Delete(db *gorm.DB, id uint) error {
	result := db.Delete(&entities.Price{}, id) // Deletes the route by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no route found to delete") // Custom error for non-existent ID
	}
	return nil
}
