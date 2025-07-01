package repository

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.User{}).Count(&total)
	return total, result.Error
}

func (r *UserRepository) Insert(ctx context.Context, conn gotann.Connection, user *domain.User) error {
	result := conn.Create(user)
	return result.Error
}

func (r *UserRepository) InsertBulk(ctx context.Context, conn gotann.Connection, users []*domain.User) error {
	result := conn.Create(&users)
	return result.Error
}

func (r *UserRepository) Update(ctx context.Context, conn gotann.Connection, user *domain.User) error {
	result := conn.Save(user)
	return result.Error
}

func (r *UserRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, users []*domain.User) error {
	result := conn.Save(&users)
	return result.Error
}

func (r *UserRepository) Delete(ctx context.Context, conn gotann.Connection, user *domain.User) error {
	result := conn.Select(clause.Associations).Delete(user)
	return result.Error
}

func (r *UserRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.User, error) {
	var users []*domain.User
	query := conn.Model(&domain.User{}).Preload("Role")
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

func (r *UserRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.User, error) {
	user := new(domain.User)
	result := conn.Preload("Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, conn gotann.Connection, email string) (*domain.User, error) {
	user := new(domain.User)
	result := conn.Preload("Role").Where("email = ? ", email).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, conn gotann.Connection, username string) (*domain.User, error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}
	fmt.Printf("username: %q\n", username)

	db := conn.Unwrap()
	if db == nil {
		return nil, errors.New("unwrapped gorm.DB is nil")
	}

	user := new(domain.User)
	result := db.Preload("Role").Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	fmt.Printf("[FindByUsername] result.RowsAffected = %d, user = %+v\n", result.RowsAffected, user)

	return user, nil
}
