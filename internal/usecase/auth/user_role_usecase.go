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
