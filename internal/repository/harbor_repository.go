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

func (hr *HarborRepository) GetAll(db *gorm.DB) ([]*entity.Harbor, error) {
	harbors := []*entity.Harbor{}
	result := db.Find(&harbors)
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
