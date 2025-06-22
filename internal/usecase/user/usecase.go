package user

import (
	"context"
	"errors"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type UserUsecase struct {
	DB             *gorm.DB
	UserRepository UserRepository
}

func NewUserUsecase(
	db *gorm.DB,
	user_repository UserRepository,
) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		UserRepository: user_repository,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, request *model.WriteUserRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		RoleID:   request.RoleID,
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		FullName: request.FullName,
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (u *UserUsecase) GetAllUsers(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadUserResponse, int, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := u.UserRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	users, err := u.UserRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all users: %w", err)
	}

	responses := make([]*model.ReadUserResponse, len(users))
	for i, user := range users {
		responses[i] = UserToResponse(user)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id uint) (*model.ReadUserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	user, err := u.UserRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return UserToResponse(user), nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, request *model.UpdateUserRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	user, err := u.UserRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Username = request.Username
	user.Email = request.Email
	user.Password = request.Password
	user.FullName = request.FullName
	user.RoleID = request.RoleID

	if err := u.UserRepository.Update(tx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	user, err := u.UserRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if err := u.UserRepository.Delete(tx, user); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
