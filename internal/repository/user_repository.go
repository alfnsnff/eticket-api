package repository

import (
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.User{}).Count(&total)
	return total, result.Error
}

func (ar *UserRepository) Insert(db *gorm.DB, user *domain.User) error {
	result := db.Create(user)
	return result.Error
}

func (ar *UserRepository) InsertBulk(db *gorm.DB, users []*domain.User) error {
	result := db.Create(&users)
	return result.Error
}

func (ur *UserRepository) Update(db *gorm.DB, user *domain.User) error {
	result := db.Save(user)
	return result.Error
}

func (ar *UserRepository) UpdateBulk(db *gorm.DB, users []*domain.User) error {
	result := db.Save(&users)
	return result.Error
}

func (ur *UserRepository) Delete(db *gorm.DB, user *domain.User) error {
	result := db.Select(clause.Associations).Delete(user)
	return result.Error
}

func (ur *UserRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.User, error) {
	var users []*domain.User
	query := db.Preload("Role")
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", search, search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (ur *UserRepository) FindByID(db *gorm.DB, id uint) (*domain.User, error) {
	user := new(domain.User)
	result := db.Preload("Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) FindByEmail(db *gorm.DB, email string) (*domain.User, error) {
	user := new(domain.User)
	result := db.Preload("Role").Where("email = ? ", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) FindByUsername(db *gorm.DB, username string) (*domain.User, error) {
	user := new(domain.User)
	result := db.Preload("Role").Where("username = ? ", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
