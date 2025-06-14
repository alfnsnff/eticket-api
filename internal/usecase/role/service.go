package role

import (
	"context"
	"errors"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type RoleUsecase struct {
	Tx             *tx.TxManager
	RoleRepository *repository.RoleRepository
}

func NewRoleUsecase(
	tx *tx.TxManager,
	role_repository *repository.RoleRepository,
) *RoleUsecase {
	return &RoleUsecase{
		Tx:             tx,
		RoleRepository: role_repository,
	}
}

func (r *RoleUsecase) CreateRole(ctx context.Context, request *model.WriteRoleRequest) error {
	user := mapper.RoleMapper.FromWrite(request)

	if user.RoleName == "" {
		return fmt.Errorf("role name cannot be empty")
	}

	if user.Description == "" {
		return fmt.Errorf("desription cannot be empty")
	}

	return r.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return r.RoleRepository.Create(tx, user)
	})
}

func (r *RoleUsecase) GetAllRoles(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadRoleResponse, int, error) {
	roles := []*entity.Role{}
	var total int64
	err := r.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = r.RoleRepository.Count(tx)
		if err != nil {
			return err
		}
		roles, err = r.RoleRepository.GetAll(tx, limit, offset, sort, search)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all allocations: %w", err)
	}

	return mapper.RoleMapper.ToModels(roles), int(total), nil
}

func (r *RoleUsecase) GetRoleByID(ctx context.Context, id uint) (*model.ReadRoleResponse, error) {
	role := new(entity.Role)

	err := r.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

func (r *RoleUsecase) UpdateRole(ctx context.Context, id uint, request *model.UpdateRoleRequest) error {
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

	return r.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return r.RoleRepository.Update(tx, role)
	})

}

func (r *RoleUsecase) DeleteRole(ctx context.Context, id uint) error {

	return r.Tx.Execute(ctx, func(tx *gorm.DB) error {
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
