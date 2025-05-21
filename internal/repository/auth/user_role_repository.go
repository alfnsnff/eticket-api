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

func (urr *UserRoleRepository) Count(db *gorm.DB) (int64, error) {
	userRoles := []*entity.UserRole{}
	var total int64
	result := db.Find(&userRoles).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (urr *UserRoleRepository) GetAll(db *gorm.DB) ([]*entity.UserRole, error) {
	user_roles := []*entity.UserRole{}
	result := db.Find(&user_roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return user_roles, nil
}

func (urr *UserRoleRepository) GetByID(db *gorm.DB, id uint) (*entity.UserRole, error) {
	user_role := new(entity.UserRole)
	result := db.Find(&user_role)
	if result.Error != nil {
		return nil, result.Error
	}
	return user_role, nil
}
