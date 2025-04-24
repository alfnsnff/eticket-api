package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type PriceRepository struct {
	DB *gorm.DB
}

func NewPriceRepository(db *gorm.DB) entities.PriceRepositoryInterface {
	return &PriceRepository{DB: db}
}

// Create inserts a new route into the database
func (p *PriceRepository) Create(price *entities.Price) error {
	result := p.DB.Create(price)
	return result.Error
}

// GetAll retrieves all routes from the database
func (p *PriceRepository) GetAll() ([]*entities.Price, error) {
	var Prices []*entities.Price
	result := p.DB.Find(&Prices) // Corrected Preload
	if result.Error != nil {
		return nil, result.Error
	}
	return Prices, nil
}

// GetByID retrieves a route by its ID
func (p *PriceRepository) GetByID(id uint) (*entities.Price, error) {
	var Price entities.Price
	result := p.DB.First(&Price, id) // Fetches the route by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return &Price, result.Error
}

// GetByID retrieves a route by its ID
func (p *PriceRepository) GetByIDs(ids []uint) ([]*entities.Price, error) {
	var Prices []*entities.Price

	result := p.DB.Where("id IN ?", ids).Preload("ShipClass").
		Preload("ShipClass.Class").Find(&Prices) // Fetches the route by ID

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return Prices, result.Error
}

// PriceRepository.go
func (r *PriceRepository) GetByRouteID(routeID uint) ([]*entities.Price, error) {
	var Prices []*entities.Price
	result := r.DB.
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
func (p *PriceRepository) Update(route *entities.Price) error {
	// Uses Gorm's Save method to update the Route
	result := p.DB.Save(route)
	return result.Error
}

// Delete removes a route from the database by its ID
func (p *PriceRepository) Delete(id uint) error {
	result := p.DB.Delete(&entities.Price{}, id) // Deletes the route by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no route found to delete") // Custom error for non-existent ID
	}
	return nil
}
