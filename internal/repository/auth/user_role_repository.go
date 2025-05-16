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
