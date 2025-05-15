package repository

import (
	entity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	repository.Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (rr *UserRepository) GetAll(db *gorm.DB) ([]*entity.User, error) {
	users := []*entity.User{}
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
