package tests

import (
	"context"
	"errors"
	"testing"

	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/domain/mocks"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"

	"github.com/glebarez/sqlite"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupUserTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite: " + err.Error())
	}
	db.AutoMigrate(&domain.User{})
	return db
}

func TestUserUsecase_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)
	req := &model.WriteUserRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashedpassword",
		FullName: "John Doe",
		RoleID:   1,
	}
	mockUserRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, user *domain.User) error {
			assert.Equal(t, "johndoe", user.Username)
			assert.Equal(t, "john@example.com", user.Email)
			assert.NotEqual(t, "hashedpassword", user.Password)
			assert.True(t, utils.CheckPasswordHash("hashedpassword", user.Password))
			assert.Equal(t, "John Doe", user.FullName)
			assert.Equal(t, uint(1), user.RoleID)
			return nil
		}).Times(1)

	err := usecase.CreateUser(context.Background(), req)

	assert.NoError(t, err)
}

func TestUserUsecase_CreateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)

	req := &model.WriteUserRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashedpassword",
		FullName: "John Doe",
		RoleID:   1,
	}

	mockUserRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(errors.New("database error")).Times(1)

	err := usecase.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestUserUsecase_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)

	expectedUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		FullName: "John Doe",
		RoleID:   1,
		Role: domain.Role{
			ID:       1,
			RoleName: "Admin",
		},
	}

	mockUserRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(expectedUser, nil).
		Times(1)

	result, err := usecase.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "johndoe", result.Username)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, "John Doe", result.FullName)
}

func TestUserUsecase_GetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)

	mockUserRepo.EXPECT().
		FindByID(gomock.Any(), uint(999)).
		Return(nil, nil).
		Times(1)

	result, err := usecase.GetUserByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestUserUsecase_ListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)

	mockUsers := []*domain.User{
		{
			ID:       1,
			Username: "johndoe",
			Email:    "john@example.com",
			FullName: "John Doe",
			RoleID:   1,
		},
		{
			ID:       2,
			Username: "janedoe",
			Email:    "jane@example.com",
			FullName: "Jane Doe",
			RoleID:   2,
		},
	}

	mockUserRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(2), nil).
		Times(1)

	mockUserRepo.EXPECT().
		FindAll(gomock.Any(), 10, 0, "id", "").
		Return(mockUsers, nil).
		Times(1)

	results, count, err := usecase.ListUsers(context.Background(), 10, 0, "id", "")

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)
	assert.Equal(t, "johndoe", results[0].Username)
	assert.Equal(t, "janedoe", results[1].Username)
}

func TestUserUsecase_UpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)

	updateRequest := &model.UpdateUserRequest{
		ID:       1,
		Username: "johndoe_updated",
		Email:    "john.updated@example.com",
		Password: "newhashedpassword",
		FullName: "John Doe Updated",
		RoleID:   2,
	}

	existingUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "oldhashedpassword",
		FullName: "John Doe",
		RoleID:   1,
	}

	mockUserRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingUser, nil).
		Times(1)

	mockUserRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, user *domain.User) error {
			assert.Equal(t, uint(1), user.ID)
			assert.Equal(t, "johndoe_updated", user.Username)
			assert.Equal(t, "john.updated@example.com", user.Email)
			assert.Equal(t, "newhashedpassword", user.Password)
			assert.Equal(t, "John Doe Updated", user.FullName)
			assert.Equal(t, uint(2), user.RoleID)
			return nil
		}).
		Times(1)

	err := usecase.UpdateUser(context.Background(), updateRequest)

	assert.NoError(t, err)
}

func TestUserUsecase_DeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	db := setupUserTestDatabase()

	usecase := usecase.NewUserUsecase(db, mockUserRepo)

	existingUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		FullName: "John Doe",
		RoleID:   1,
	}

	mockUserRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingUser, nil).
		Times(1)

	mockUserRepo.EXPECT().
		Delete(gomock.Any(), existingUser).
		Return(nil).
		Times(1)

	err := usecase.DeleteUser(context.Background(), 1)

	assert.NoError(t, err)
}
