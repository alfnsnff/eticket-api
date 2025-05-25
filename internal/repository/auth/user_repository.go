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

func (ur *UserRepository) Count(db *gorm.DB) (int64, error) {
	users := []*entity.User{}
	var total int64
	result := db.Find(&users).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (ur *UserRepository) GetAll(db *gorm.DB) ([]*entity.User, error) {
	users := []*entity.User{}
	result := db.Preload("Role").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (ur *UserRepository) GetByUsername(db *gorm.DB, username string) (*entity.User, error) {
	user := new(entity.User)
	result := db.Preload("Role").Where("username = ? ", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) GetByID(db *gorm.DB, id uint) (*entity.User, error) {
	user := new(entity.User)
	result := db.Preload("Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
