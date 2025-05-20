package usecase

import (
	"context"
	"errors"
	authentity "eticket-api/internal/domain/entity/auth"
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

func (ur *UserUsecase) GetAllUsers(ctx context.Context) ([]*authmodel.ReadUserResponse, error) {
	users := []*authentity.User{}

	err := tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		var err error
		users, err = ur.UserRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, err
	}

	if users == nil {
		return nil, errors.New("user role not found")
	}

	return mapper.UserMapper.ToModels(users), nil
}

func (ur *UserUsecase) GetUserByID(ctx context.Context, id uint) (*authmodel.ReadUserResponse, error) {
	user := new(authentity.User)

	err := tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		var err error
		user, err = ur.UserRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user role not found")
	}

	return mapper.UserMapper.ToModel(user), nil
}

func (ur *UserUsecase) UpdateUser(ctx context.Context, id uint, request *authmodel.UpdateUserRequest) error {
	user := mapper.UserMapper.FromUpdate(request)
	user.ID = id

	if user.ID == 0 {
		return fmt.Errorf("role ID cannot be zero")
	}

	if user.Username == "" {
		return fmt.Errorf("role name cannot be empty")
	}

	if user.Email == "" {
		return fmt.Errorf("desription cannot be empty")
	}

	if user.Password == "" {
		return fmt.Errorf("desription cannot be empty")
	}

	return tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		return ur.UserRepository.Update(tx, user)
	})
}

func (ur *UserUsecase) DeleteUser(ctx context.Context, id uint) error {

	return tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		role, err := ur.UserRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.New("route not found")
		}
		return ur.UserRepository.Delete(tx, role)
	})

}
