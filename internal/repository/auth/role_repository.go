package repository

import (
	entity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type RoleRepository struct {
	repository.Repository[entity.Role]
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (rr *RoleRepository) GetAll(db *gorm.DB) ([]*entity.Role, error) {
	roles := []*entity.Role{}
	result := db.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (rr *RoleRepository) GetByID(db *gorm.DB, id uint) (*entity.Role, error) {
	role := new(entity.Role)
	result := db.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}
