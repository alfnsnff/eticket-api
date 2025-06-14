package repository

import (
	"errors"
	"strings"

	"eticket-api/internal/entity"

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

func (hr *HarborRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Harbor, error) {
	harbors := []*entity.Harbor{}

	query := db

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("harbor_name ILIKE ? OR harbor_alias ILIKE ?", search, search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&harbors).Error
	return harbors, err
}

func (hr *HarborRepository) GetByID(db *gorm.DB, id uint) (*entity.Harbor, error) {
	harbor := new(entity.Harbor)
	result := db.First(&harbor, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return harbor, result.Error
}
