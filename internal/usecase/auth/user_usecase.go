package usecase

import (
	"context"
	"errors"
	authentity "eticket-api/internal/domain/entity/auth"
	authmodel "eticket-api/internal/model/auth"
	"eticket-api/internal/model/mapper"
	authrepository "eticket-api/internal/repository/auth"
	utils "eticket-api/pkg/utils/hash"
	"eticket-api/pkg/utils/tx"
	"fmt"

	"gorm.io/gorm"
)

type UserUsecase struct {
	Tx             *tx.TxManager
	UserRepository *authrepository.UserRepository
}

func NewUserUsecase(
	tx *tx.TxManager,
	user_repository *authrepository.UserRepository,
) *UserUsecase {
	return &UserUsecase{
		Tx:             tx,
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

	return u.Tx.Execute(ctx, func(tx *gorm.DB) error {
		hashedPassword, err := utils.HashPassword(request.Password) // Use the helper from utils
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
		return u.UserRepository.Create(tx, user)
	})
}

func (u *UserUsecase) GetAllUsers(ctx context.Context, limit, offset int) ([]*authmodel.ReadUserResponse, int, error) {
	users := []*authentity.User{}
	var total int64
	err := u.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = u.UserRepository.Count(tx)
		if err != nil {
			return err
		}
		users, err = u.UserRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return mapper.UserMapper.ToModels(users), int(total), nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id uint) (*authmodel.ReadUserResponse, error) {
	user := new(authentity.User)

	err := u.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		user, err = u.UserRepository.GetByID(tx, id)
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

func (u *UserUsecase) UpdateUser(ctx context.Context, id uint, request *authmodel.UpdateUserRequest) error {
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

	return u.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return u.UserRepository.Update(tx, user)
	})
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uint) error {

	return u.Tx.Execute(ctx, func(tx *gorm.DB) error {
		role, err := u.UserRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.New("route not found")
		}
		return u.UserRepository.Delete(tx, role)
	})

}
