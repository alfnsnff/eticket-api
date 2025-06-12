package repository

import (
	"eticket-api/internal/entity"

	"gorm.io/gorm"
)

type RoleRepository struct {
	Repository[entity.Role]
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
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

func (ror *RoleRepository) GetAll(db *gorm.DB) ([]*entity.Role, error) {
	roles := []*entity.Role{}
	result := db.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (ror *RoleRepository) GetByID(db *gorm.DB, id uint) (*entity.Role, error) {
	role := new(entity.Role)
	result := db.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}
