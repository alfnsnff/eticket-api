package repository

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShipRepository struct {
	DB *gorm.DB
}

func NewShipRepository(db *gorm.DB) *ShipRepository {
	return &ShipRepository{DB: db}
}

func (r *ShipRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Ship{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *ShipRepository) Insert(ctx context.Context, conn gotann.Connection, ship *domain.Ship) error {
	result := conn.Create(ship)
	return result.Error
}

func (r *ShipRepository) InsertBulk(ctx context.Context, conn gotann.Connection, ships []*domain.Ship) error {
	result := conn.Create(ships)
	return result.Error
}

func (r *ShipRepository) Update(ctx context.Context, conn gotann.Connection, ship *domain.Ship) error {
	result := conn.Save(ship)
	return result.Error
}

func (r *ShipRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, ships []*domain.Ship) error {
	result := conn.Save(ships)
	return result.Error
}

func (r *ShipRepository) Delete(ctx context.Context, conn gotann.Connection, ship *domain.Ship) error {
	result := conn.Select(clause.Associations).Delete(ship)
	return result.Error
}

func (r *ShipRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Ship, error) {
	ships := []*domain.Ship{}
	query := conn.Model(&domain.Ship{})
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

func (r *ShipRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Ship, error) {
	ship := new(domain.Ship)
	result := conn.First(&ship, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ship, result.Error
}
