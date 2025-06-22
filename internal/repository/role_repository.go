package repository

import (
	"eticket-api/internal/entity"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (ar *RoleRepository) Create(db *gorm.DB, role *entity.Role) error {
	result := db.Create(role)
	return result.Error
}

func (ar *RoleRepository) Update(db *gorm.DB, role *entity.Role) error {
	result := db.Save(role)
	return result.Error
}

func (ar *RoleRepository) Delete(db *gorm.DB, role *entity.Role) error {
	result := db.Select(clause.Associations).Delete(role)
	return result.Error
}

func (ror *RoleRepository) Count(db *gorm.DB) (int64, error) {
	roles := []*entity.Role{}
	var total int64
	result := db.Find(&roles).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (ror *RoleRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Role, error) {
	roles := []*entity.Role{}

	query := db

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("role_name ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&roles).Error
	return roles, err
}
func (ror *RoleRepository) GetByID(db *gorm.DB, id uint) (*entity.Role, error) {
	role := new(entity.Role)
	result := db.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}
