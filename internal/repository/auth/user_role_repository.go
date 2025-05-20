package repository

import (
	entity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type UserRoleRepository struct {
	repository.Repository[entity.UserRole]
}

func NewUserRoleRepository() *UserRoleRepository {
	return &UserRoleRepository{}
}

func (rr *UserRoleRepository) GetAll(db *gorm.DB) ([]*entity.UserRole, error) {
	user_roles := []*entity.UserRole{}
	result := db.Find(&user_roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return user_roles, nil
}

func (rr *UserRoleRepository) GetByID(db *gorm.DB, id uint) (*entity.UserRole, error) {
	user_role := new(entity.UserRole)
	result := db.Find(&user_role)
	if result.Error != nil {
		return nil, result.Error
	}
	return user_role, nil
}
