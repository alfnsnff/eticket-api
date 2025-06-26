package repository

import (
	"errors"
	"strings"

	"eticket-api/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HarborRepository struct{}

func NewHarborRepository() *HarborRepository {
	return &HarborRepository{}
}

func (hr *HarborRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Harbor{}).Count(&total)
	return total, result.Error
}

func (hr *HarborRepository) Insert(db *gorm.DB, harbor *domain.Harbor) error {
	result := db.Create(harbor)
	return result.Error
}

func (hr *HarborRepository) InsertBulk(db *gorm.DB, harbors []*domain.Harbor) error {
	result := db.Create(harbors)
	return result.Error
}

func (hr *HarborRepository) Update(db *gorm.DB, harbor *domain.Harbor) error {
	result := db.Save(harbor)
	return result.Error
}

func (hr *HarborRepository) UpdateBulk(db *gorm.DB, harbors []*domain.Harbor) error {
	result := db.Save(harbors)
	return result.Error
}

func (hr *HarborRepository) Delete(db *gorm.DB, harbor *domain.Harbor) error {
	result := db.Select(clause.Associations).Delete(harbor)
	return result.Error
}

func (hr *HarborRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Harbor, error) {
	harbors := []*domain.Harbor{}
	query := db
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("harbor_name ILIKE ? OR harbor_alias ILIKE ?", search, search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&harbors).Error
	return harbors, err
}

func (hr *HarborRepository) FindByID(db *gorm.DB, id uint) (*domain.Harbor, error) {
	harbor := new(domain.Harbor)
	result := db.First(&harbor, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return harbor, result.Error
}
