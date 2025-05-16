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

type RoleUsecase struct {
	DB             *gorm.DB
	RoleRepository *authrepository.RoleRepository
}

func NewRoleUsecase(
	db *gorm.DB,
	role_repository *authrepository.RoleRepository,
) *RoleUsecase {
	return &RoleUsecase{
		DB:             db,
		RoleRepository: role_repository,
	}
}

func (r *RoleUsecase) CreateRole(ctx context.Context, request *authmodel.WriteRoleRequest) error {
	user := mapper.RoleMapper.FromWrite(request)

	if user.RoleName == "" {
		return fmt.Errorf("role name cannot be empty")
	}

	if user.Description == "" {
		return fmt.Errorf("desription cannot be empty")
	}

	return tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		return r.RoleRepository.Create(tx, user)
	})
}
