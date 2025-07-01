package repository

import (
	"context"
	"errors"
	"strings"

	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HarborRepository struct {
	DB *gorm.DB
}

func NewHarborRepository(db *gorm.DB) *HarborRepository {
	return &HarborRepository{DB: db}
}

func (r *HarborRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Harbor{}).Count(&total)
	return total, result.Error
}

func (r *HarborRepository) Insert(ctx context.Context, conn gotann.Connection, harbor *domain.Harbor) error {
	result := conn.Create(harbor)
	return result.Error
}

func (r *HarborRepository) InsertBulk(ctx context.Context, conn gotann.Connection, harbors []*domain.Harbor) error {
	result := conn.Create(harbors)
	return result.Error
}

func (r *HarborRepository) Update(ctx context.Context, conn gotann.Connection, harbor *domain.Harbor) error {
	result := conn.Save(harbor)
	return result.Error
}

func (r *HarborRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, harbors []*domain.Harbor) error {
	result := conn.Save(harbors)
	return result.Error
}

func (r *HarborRepository) Delete(ctx context.Context, conn gotann.Connection, harbor *domain.Harbor) error {
	result := conn.Select(clause.Associations).Delete(harbor)
	return result.Error
}

func (r *HarborRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Harbor, error) {
	harbors := []*domain.Harbor{}
	query := conn.Model(&domain.Harbor{})
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

func (r *HarborRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Harbor, error) {
	harbor := new(domain.Harbor)
	result := conn.First(&harbor, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return harbor, result.Error
}
