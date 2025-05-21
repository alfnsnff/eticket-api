package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type HarborRepository struct {
	Repository[entity.Harbor]
}

func NewHarborRepository() *HarborRepository {
	return &HarborRepository{}
}

func (hr *HarborRepository) Count(db *gorm.DB) (int64, error) {
	harbors := []*entity.Harbor{}
	var total int64
	result := db.Find(&harbors).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (hr *HarborRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Harbor, error) {
	harbors := []*entity.Harbor{}
	result := db.Limit(limit).Offset(offset).Find(&harbors)
	if result.Error != nil {
		return nil, result.Error
	}
	return harbors, nil
}

func (hr *HarborRepository) GetByID(db *gorm.DB, id uint) (*entity.Harbor, error) {
	harbor := new(entity.Harbor)
	result := db.First(&harbor, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return harbor, result.Error
}
