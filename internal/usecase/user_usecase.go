package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type UserUsecase struct {
	Transactor     *transact.Transactor
	UserRepository domain.UserRepository
}

func NewUserUsecase(

	transactor *transact.Transactor, // Assuming transact package is imported
	user_repository domain.UserRepository,
) *UserUsecase {
	return &UserUsecase{

		Transactor:     transactor,
		UserRepository: user_repository,
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, e *domain.User) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		hashedPassword, err := utils.HashPassword(e.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user := &domain.User{
			RoleID:   e.RoleID,
			Username: e.Username,
			Email:    e.Email,
			Password: hashedPassword,
			FullName: e.FullName,
		}

		if err := uc.UserRepository.Insert(ctx, tx, user); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
}

func (uc *UserUsecase) ListUsers(ctx context.Context, limit, offset int, sort, search string) ([]*domain.User, int, error) {
	var err error
	var total int64
	var users []*domain.User
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.UserRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count users: %w", err)
		}

		users, err = uc.UserRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all users: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	return users, int(total), nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {

	var err error
	var user *domain.User
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		fmt.Printf("[Login] tx type: %T, value: %#v\n", tx, tx)
		user, err = uc.UserRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, e *domain.User) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		user, err := uc.UserRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if user == nil {
			return errs.ErrNotFound
		}

		user.Username = e.Username
		user.Email = e.Email
		user.Password = e.Password
		user.FullName = e.FullName
		user.RoleID = e.RoleID

		if err := uc.UserRepository.Update(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
		return nil
	})
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		user, err := uc.UserRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return errs.ErrNotFound
		}

		if err := uc.UserRepository.Delete(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return nil
	})
}
