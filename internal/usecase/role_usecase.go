package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type RoleUsecase struct {
	DB             *gorm.DB
	RoleRepository domain.RoleRepository
}

func NewRoleUsecase(
	db *gorm.DB,
	roleRepository domain.RoleRepository,
) *RoleUsecase {
	return &RoleUsecase{
		DB:             db,
		RoleRepository: roleRepository,
	}
}

func (r *RoleUsecase) CreateRole(ctx context.Context, request *model.WriteRoleRequest) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	role := &domain.Role{
		RoleName:    request.RoleName,
		Description: request.Description,
	}

	if err := r.RoleRepository.Insert(tx, role); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *RoleUsecase) ListRoles(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadRoleResponse, int, error) {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := r.RoleRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	roles, err := r.RoleRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all roles: %w", err)
	}

	responses := make([]*model.ReadRoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = mapper.RoleToResponse(role)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit: %w", err)
	}

	return responses, int(total), nil
}

func (r *RoleUsecase) GetRoleByID(ctx context.Context, id uint) (*model.ReadRoleResponse, error) {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	role, err := r.RoleRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by id: %w", err)
	}
	if role == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return mapper.RoleToResponse(role), nil
}

func (r *RoleUsecase) UpdateRole(ctx context.Context, request *model.UpdateRoleRequest) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Fetch existing allocation
	role, err := r.RoleRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}
	if role == nil {
		return errs.ErrNotFound
	}

	role.RoleName = request.RoleName
	role.Description = request.Description

	if err := r.RoleRepository.Update(tx, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *RoleUsecase) DeleteRole(ctx context.Context, id uint) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	role, err := r.RoleRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errs.ErrNotFound
	}

	if err := r.RoleRepository.Delete(tx, role); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
