package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type RoleUsecase struct {
	Transactor     transact.Transactor // Uncomment if you need transaction management
	RoleRepository domain.RoleRepository
}

func NewRoleUsecase(
	transactor transact.Transactor, // Uncomment if you need transaction management
	roleRepository domain.RoleRepository,
) *RoleUsecase {
	return &RoleUsecase{
		Transactor:     transactor, // Uncomment if you need transaction management
		RoleRepository: roleRepository,
	}
}

func (uc *RoleUsecase) CreateRole(ctx context.Context, e *domain.Role) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		role := &domain.Role{
			RoleName:    e.RoleName,
			Description: e.Description,
		}

		if err := uc.RoleRepository.Insert(ctx, tx, role); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create role: %w", err)
		}
		return nil
	})
}

func (uc *RoleUsecase) ListRoles(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Role, int, error) {
	var err error
	var total int64
	var roles []*domain.Role
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.RoleRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count roles: %w", err)
		}

		roles, err = uc.RoleRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all roles: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	return roles, int(total), nil
}

func (uc *RoleUsecase) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	var err error
	var role *domain.Role
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		role, err = uc.RoleRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get role by id: %w", err)
		}
		if role == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get role by id: %w", err)
	}
	return role, nil
}

func (uc *RoleUsecase) UpdateRole(ctx context.Context, e *domain.Role) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		role, err := uc.RoleRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find role: %w", err)
		}
		if role == nil {
			return errs.ErrNotFound
		}

		role.RoleName = e.RoleName
		role.Description = e.Description

		if err := uc.RoleRepository.Update(ctx, tx, role); err != nil {
			return fmt.Errorf("failed to update role: %w", err)
		}

		return nil
	})
}

func (uc *RoleUsecase) DeleteRole(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		role, err := uc.RoleRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get role: %w", err)
		}
		if role == nil {
			return errs.ErrNotFound
		}

		if err := uc.RoleRepository.Delete(ctx, tx, role); err != nil {
			return fmt.Errorf("failed to delete role: %w", err)
		}
		return nil
	})
}
