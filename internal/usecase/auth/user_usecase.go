package usecase

import (
	"context"
	authmodel "eticket-api/internal/model/auth"
	"eticket-api/internal/model/mapper"
	authrepository "eticket-api/internal/repository/auth"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type UserUsecase struct {
	DB             *gorm.DB
	UserRepository *authrepository.UserRepository
}

func NewUserUsecase(
	db *gorm.DB,
	user_repository *authrepository.UserRepository,
) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		UserRepository: user_repository,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, request *authmodel.WriteUserRequest) error {
	user := mapper.UserMapper.FromWrite(request)

	if user.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if user.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if user.FullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}

	return tx.Execute(ctx, u.DB, func(tx *gorm.DB) error {
		return u.UserRepository.Create(tx, user)
	})
}
