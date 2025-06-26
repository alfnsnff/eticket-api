package tests

import (
	"context"
	"errors"
	"testing"

	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/domain/mocks"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"

	"github.com/glebarez/sqlite"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRoleTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite: " + err.Error())
	}

	// Auto migrate the schema
	db.AutoMigrate(&domain.Role{})

	return db
}

func TestRoleUsecase_CreateRole_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	req := &model.WriteRoleRequest{
		RoleName:    "admin",
		Description: "Administrator role",
	}

	mockRoleRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, role *domain.Role) error {
			assert.Equal(t, "admin", role.RoleName)
			assert.Equal(t, "Administrator role", role.Description)
			role.ID = 1 // Simulate database auto-assignment
			return nil
		}).Times(1)

	err := usecase.CreateRole(context.Background(), req)

	assert.NoError(t, err)
}

func TestRoleUsecase_CreateRole_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	req := &model.WriteRoleRequest{
		RoleName:    "admin",
		Description: "Administrator role",
	}

	mockRoleRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(errors.New("database error")).Times(1)

	err := usecase.CreateRole(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create role")
}

func TestRoleUsecase_ListRoles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	roles := []*domain.Role{
		{
			ID:          1,
			RoleName:    "admin",
			Description: "Administrator role",
		},
		{
			ID:          2,
			RoleName:    "user",
			Description: "Regular user role",
		},
	}

	mockRoleRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(2), nil).Times(1)

	mockRoleRepo.EXPECT().
		FindAll(gomock.Any(), 10, 0, "id asc", "").
		Return(roles, nil).Times(1)

	responses, total, err := usecase.ListRoles(context.Background(), 10, 0, "id asc", "")

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, responses, 2)
	assert.Equal(t, uint(1), responses[0].ID)
	assert.Equal(t, "admin", responses[0].RoleName)
	assert.Equal(t, "Administrator role", responses[0].Description)
}

func TestRoleUsecase_ListRoles_CountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	mockRoleRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(0), errors.New("count error")).Times(1)

	responses, total, err := usecase.ListRoles(context.Background(), 10, 0, "id asc", "")

	assert.Error(t, err)
	assert.Nil(t, responses)
	assert.Equal(t, 0, total)
	assert.Contains(t, err.Error(), "failed to count roles")
}

func TestRoleUsecase_GetRoleByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	role := &domain.Role{
		ID:          1,
		RoleName:    "admin",
		Description: "Administrator role",
	}

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(role, nil).Times(1)

	response, err := usecase.GetRoleByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "admin", response.RoleName)
	assert.Equal(t, "Administrator role", response.Description)
}

func TestRoleUsecase_GetRoleByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), uint(999)).
		Return(nil, nil).Times(1)

	response, err := usecase.GetRoleByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestRoleUsecase_GetRoleByID_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(nil, errors.New("database error")).Times(1)

	response, err := usecase.GetRoleByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get role by id")
}

func TestRoleUsecase_UpdateRole_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	req := &model.UpdateRoleRequest{
		ID:          1,
		RoleName:    "updated_admin",
		Description: "Updated administrator role",
	}

	existingRole := &domain.Role{
		ID:          1,
		RoleName:    "admin",
		Description: "Administrator role",
	}

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), req.ID).
		Return(existingRole, nil).Times(1)

	mockRoleRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, role *domain.Role) error {
			assert.Equal(t, uint(1), role.ID)
			assert.Equal(t, "updated_admin", role.RoleName)
			assert.Equal(t, "Updated administrator role", role.Description)
			return nil
		}).Times(1)

	err := usecase.UpdateRole(context.Background(), req)

	assert.NoError(t, err)
}

func TestRoleUsecase_UpdateRole_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	req := &model.UpdateRoleRequest{
		ID:          999,
		RoleName:    "updated_admin",
		Description: "Updated administrator role",
	}

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), req.ID).
		Return(nil, nil).Times(1)

	err := usecase.UpdateRole(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestRoleUsecase_UpdateRole_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	req := &model.UpdateRoleRequest{
		ID:          1,
		RoleName:    "updated_admin",
		Description: "Updated administrator role",
	}

	existingRole := &domain.Role{
		ID:          1,
		RoleName:    "admin",
		Description: "Administrator role",
	}

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), req.ID).
		Return(existingRole, nil).Times(1)

	mockRoleRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errors.New("update error")).Times(1)

	err := usecase.UpdateRole(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update role")
}

func TestRoleUsecase_DeleteRole_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	role := &domain.Role{
		ID:          1,
		RoleName:    "admin",
		Description: "Administrator role",
	}

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(role, nil).Times(1)

	mockRoleRepo.EXPECT().
		Delete(gomock.Any(), role).
		Return(nil).Times(1)

	err := usecase.DeleteRole(context.Background(), 1)

	assert.NoError(t, err)
}

func TestRoleUsecase_DeleteRole_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	db := setupRoleTestDatabase()

	usecase := usecase.NewRoleUsecase(db, mockRoleRepo)

	mockRoleRepo.EXPECT().
		FindByID(gomock.Any(), uint(999)).
		Return(nil, nil).Times(1)

	err := usecase.DeleteRole(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}
