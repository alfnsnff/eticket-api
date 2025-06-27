package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClaimItemRepository struct{}

func NewClaimItemRepository() *ClaimItemRepository {
	return &ClaimItemRepository{}
}

func (cr *ClaimItemRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.ClaimItem{}).Count(&total)
	return total, result.Error
}

func (ar *ClaimItemRepository) Insert(db *gorm.DB, claimItem *domain.ClaimItem) error {
	result := db.Create(claimItem)
	return result.Error
}

func (cr *ClaimItemRepository) InsertBulk(db *gorm.DB, ClaimItems []*domain.ClaimItem) error {
	result := db.Create(&ClaimItems)
	return result.Error
}

func (cr *ClaimItemRepository) Update(db *gorm.DB, claimItem *domain.ClaimItem) error {
	result := db.Save(claimItem)
	return result.Error
}

func (cr *ClaimItemRepository) UpdateBulk(db *gorm.DB, ClaimItems []*domain.ClaimItem) error {
	result := db.Save(&ClaimItems)
	return result.Error
}

func (cr *ClaimItemRepository) Delete(db *gorm.DB, claimItem *domain.ClaimItem) error {
	result := db.Select(clause.Associations).Delete(claimItem)
	return result.Error
}

func (cr *ClaimItemRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.ClaimItem, error) {
	ClaimItemes := []*domain.ClaimItem{}
	query := db
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

func (cr *ClaimItemRepository) FindByID(db *gorm.DB, id uint) (*domain.ClaimItem, error) {
	ClaimItem := new(domain.ClaimItem)
	result := db.First(&ClaimItem, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ClaimItem, result.Error
}
