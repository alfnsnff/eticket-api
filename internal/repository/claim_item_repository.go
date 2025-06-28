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

type ClaimItemRepository struct {
	DB *gorm.DB
}

func NewClaimItemRepository(db *gorm.DB) *ClaimItemRepository {
	return &ClaimItemRepository{DB: db}
}

func (r *ClaimItemRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.ClaimItem{}).Count(&total)
	return total, result.Error
}

func (r *ClaimItemRepository) Insert(ctx context.Context, conn gotann.Connection, claimItem *domain.ClaimItem) error {
	result := conn.Create(claimItem)
	return result.Error
}

func (r *ClaimItemRepository) InsertBulk(ctx context.Context, conn gotann.Connection, ClaimItems []*domain.ClaimItem) error {
	result := conn.Create(&ClaimItems)
	return result.Error
}

func (r *ClaimItemRepository) Update(ctx context.Context, conn gotann.Connection, claimItem *domain.ClaimItem) error {
	result := conn.Save(claimItem)
	return result.Error
}

func (r *ClaimItemRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, ClaimItems []*domain.ClaimItem) error {
	result := conn.Save(&ClaimItems)
	return result.Error
}

func (r *ClaimItemRepository) Delete(ctx context.Context, conn gotann.Connection, claimItem *domain.ClaimItem) error {
	result := conn.Select(clause.Associations).Delete(claimItem)
	return result.Error
}

func (r *ClaimItemRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.ClaimItem, error) {
	ClaimItemes := []*domain.ClaimItem{}
	query := conn.Model(&domain.ClaimItem{})
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("ClaimItem_name ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&ClaimItemes).Error
	return ClaimItemes, err
}

func (r *ClaimItemRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.ClaimItem, error) {
	ClaimItem := new(domain.ClaimItem)
	result := conn.First(&ClaimItem, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ClaimItem, result.Error
}
