package usecase

import (
	"context"
	"errors"
	authentity "eticket-api/internal/domain/entity/auth"
	authmodel "eticket-api/internal/model/auth"
	"eticket-api/internal/model/mapper"
	authrepository "eticket-api/internal/repository/auth"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type UserRoleUsecase struct {
	DB                 *gorm.DB
	RoleRepository     *authrepository.RoleRepository
	UserRepository     *authrepository.UserRepository
	UserRoleRepository *authrepository.UserRoleRepository
}

func NewUserRoleUsecase(
	db *gorm.DB,
	role_repository *authrepository.RoleRepository,
	user_repository *authrepository.UserRepository,
	user_role_repository *authrepository.UserRoleRepository,
) *UserRoleUsecase {
	return &UserRoleUsecase{
		DB:                 db,
		RoleRepository:     role_repository,
		UserRepository:     user_repository,
		UserRoleRepository: user_role_repository,
	}
}

func (ur *UserRoleUsecase) CreateUserRole(ctx context.Context, request *authmodel.WriteUserRoleRequest) error {
	user_role := mapper.UserRoleMapper.FromWrite(request)

	if user_role.UserID == 0 {
		return fmt.Errorf("user ID cannot be zero")
	}

	if user_role.RoleID == 0 {
		return fmt.Errorf("role ID cannot be zero")
	}

	return tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		return ur.UserRoleRepository.Create(tx, user_role)
	})
}

func (ur *UserRoleUsecase) GetAllUserRoles(ctx context.Context) ([]*authmodel.ReadUserRoleResponse, error) {
	user_roles := []*authentity.UserRole{}

	err := tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		var err error
		user_roles, err = ur.UserRoleRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, err
	}

	if user_roles == nil {
		return nil, errors.New("user role not found")
	}

	return mapper.UserRoleMapper.ToModels(user_roles), nil
}

func (ur *UserRoleUsecase) GetUserRoleByID(ctx context.Context, id uint) (*authmodel.ReadUserRoleResponse, error) {
	user_role := new(authentity.UserRole)

	err := tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		var err error
		user_role, err = ur.UserRoleRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if user_role == nil {
		return nil, errors.New("user role not found")
	}

	return mapper.UserRoleMapper.ToModel(user_role), nil
}

func (ur *UserRoleUsecase) UpdateUserRole(ctx context.Context, id uint, request *authmodel.UpdateUserRoleRequest) error {
	user_role := mapper.UserRoleMapper.FromUpdate(request)
	user_role.ID = id

	if user_role.ID == 0 {
		return fmt.Errorf("role ID cannot be zero")
	}

	if user_role.UserID == 0 {
		return fmt.Errorf("role name cannot be empty")
	}

	if user_role.RoleID == 0 {
		return fmt.Errorf("desription cannot be empty")
	}

	return tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		return ur.UserRoleRepository.Update(tx, user_role)
	})
}

func (ur *UserRoleUsecase) DeleteUserRole(ctx context.Context, id uint) error {

	return tx.Execute(ctx, ur.DB, func(tx *gorm.DB) error {
		role, err := ur.RoleRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.New("route not found")
		}
		return ur.RoleRepository.Delete(tx, role)
	})

}
