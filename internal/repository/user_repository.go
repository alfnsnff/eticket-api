package repository

import (
	"eticket-api/internal/entity"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
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

func (ur *UserRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.User, error) {
	var users []*entity.User

	query := db.Preload("Role")

	// 🔍 Search (on name or email)
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", search, search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (ur *UserRepository) GetByUsername(db *gorm.DB, username string) (*entity.User, error) {
	user := new(entity.User)
	result := db.Preload("Role").Where("username = ? ", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) GetByEmail(db *gorm.DB, email string) (*entity.User, error) {
	user := new(entity.User)
	result := db.Preload("Role").Where("email = ? ", email).First(&user)
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

func (ur *UserRepository) UpdatePassword(db *gorm.DB, userID uint, password string) error {
	users := []*entity.User{}
	result := db.First(&users, userID).Update("password", password)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
