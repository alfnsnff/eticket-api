package repository

import (
	"context"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Role{}).Count(&total)
	return total, result.Error
}

func (r *RoleRepository) Insert(ctx context.Context, conn gotann.Connection, role *domain.Role) error {
	result := conn.Create(role)
	return result.Error
}

func (r *RoleRepository) InsertBulk(ctx context.Context, conn gotann.Connection, roles []*domain.Role) error {
	result := conn.Create(&roles)
	return result.Error
}

func (r *RoleRepository) Update(ctx context.Context, conn gotann.Connection, role *domain.Role) error {
	result := conn.Save(role)
	return result.Error
}

func (r *RoleRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, roles []*domain.Role) error {
	result := conn.Save(&roles)
	return result.Error
}

func (r *RoleRepository) Delete(ctx context.Context, conn gotann.Connection, role *domain.Role) error {
	result := conn.Select(clause.Associations).Delete(role)
	return result.Error
}

func (r *RoleRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Role, error) {
	roles := []*domain.Role{}
	query := conn.Model(&domain.Ship{})
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("role_name ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&roles).Error
	return roles, err
}
func (r *RoleRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Role, error) {
	role := new(domain.Role)
	result := conn.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}
