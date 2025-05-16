package usecase

import (
	"context"
	"errors"
	authmodel "eticket-api/internal/model/auth"
	"eticket-api/internal/model/mapper"
	authrepository "eticket-api/internal/repository/auth"
	tx "eticket-api/pkg/utils/helper"
	"eticket-api/pkg/utils/helper/auth"
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
		hashedPassword, err := auth.HashPassword(request.Password) // Use the helper from utils
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
		return u.UserRepository.Create(tx, user)
	})
}

// Login authenticates a user and returns access and refresh tokens.
func (uc *UserUsecase) Login(ctx context.Context, request *authmodel.UserLoginRequest) (string, string, error) {
	if request.Username == "" || request.Password == "" {
		return "", "", errors.New("username and password are required")
	}

	var accessToken, refreshToken string

	err := tx.Execute(ctx, uc.DB, func(txDB *gorm.DB) error {
		user, repoErr := uc.UserRepository.GetByUsername(txDB, request.Username)
		if repoErr != nil {
			return fmt.Errorf("failed to retrieve user: %w", repoErr)
		}
		if user == nil {
			return errors.New("invalid credentials")
		}

		if !auth.CheckPasswordHash(request.Password, user.Password) {
			return errors.New("invalid credentials")
		}

		var err error
		accessToken, err = auth.GenerateAccessToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate access token: %w", err)
		}

		refreshToken, err = auth.GenerateRefreshToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %w", err)
		}

		return nil // ✅ Only return error inside tx
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
