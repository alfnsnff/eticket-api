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

func (r *RoleUsecase) GetAllRoles(ctx context.Context) ([]*authmodel.ReadRoleResponse, error) {
	roles := []*authentity.Role{}

	err := tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		var err error
		roles, err = r.RoleRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, err
	}

	if roles == nil {
		return nil, errors.New("booking not found")
	}

	return mapper.RoleMapper.ToModels(roles), nil
}

func (r *RoleUsecase) GetRoleByID(ctx context.Context, id uint) (*authmodel.ReadRoleResponse, error) {
	role := new(authentity.Role)

	err := tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		var err error
		role, err = r.RoleRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("booking not found")
	}

	return mapper.RoleMapper.ToModel(role), nil
}

func (r *RoleUsecase) UpdateRole(ctx context.Context, id uint, request *authmodel.UpdateRoleRequest) error {
	role := mapper.RoleMapper.FromUpdate(request)
	role.ID = id

	if role.ID == 0 {
		return fmt.Errorf("role ID cannot be zero")
	}

	if role.RoleName == "" {
		return fmt.Errorf("role name cannot be empty")
	}

	if role.Description == "" {
		return fmt.Errorf("desription cannot be empty")
	}

	return tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		return r.RoleRepository.Update(tx, role)
	})

}

func (r *RoleUsecase) DeleteRole(ctx context.Context, id uint) error {

	return tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		role, err := r.RoleRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.New("route not found")
		}
		return r.RoleRepository.Delete(tx, role)
	})

}
