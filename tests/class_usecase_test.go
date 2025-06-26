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

func setupClassTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite: " + err.Error())
	}
	db.AutoMigrate(&domain.Class{})
	return db
}

func TestClassUsecase_CreateClass_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	req := &model.WriteClassRequest{
		ClassName:  "Economy",
		Type:       "Standard",
		ClassAlias: "ECO",
	}

	mockClassRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, class *domain.Class) error {
			assert.Equal(t, "Economy", class.ClassName)
			assert.Equal(t, "Standard", class.Type)
			assert.Equal(t, "ECO", class.ClassAlias)
			return nil
		}).Times(1)

	err := usecase.CreateClass(context.Background(), req)

	assert.NoError(t, err)
}

func TestClassUsecase_CreateClass_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	req := &model.WriteClassRequest{
		ClassName:  "Economy",
		Type:       "Standard",
		ClassAlias: "ECO",
	}

	mockClassRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(errors.New("database error")).Times(1)

	err := usecase.CreateClass(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create class")
}

func TestClassUsecase_GetClassByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	expectedClass := &domain.Class{
		ID:         1,
		ClassName:  "Business",
		Type:       "Premium",
		ClassAlias: "BIZ",
	}

	mockClassRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(expectedClass, nil).
		Times(1)

	result, err := usecase.GetClassByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Business", result.ClassName)
	assert.Equal(t, "Premium", result.Type)
}

func TestClassUsecase_GetClassByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	mockClassRepo.EXPECT().
		FindByID(gomock.Any(), uint(999)).
		Return(nil, nil).
		Times(1)

	result, err := usecase.GetClassByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestClassUsecase_ListClasses_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	mockClasses := []*domain.Class{
		{
			ID:         1,
			ClassName:  "Economy",
			Type:       "Standard",
			ClassAlias: "ECO",
		},
		{
			ID:         2,
			ClassName:  "Business",
			Type:       "Premium",
			ClassAlias: "BIZ",
		},
	}

	mockClassRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(2), nil).
		Times(1)

	mockClassRepo.EXPECT().
		FindAll(gomock.Any(), 10, 0, "id", "").
		Return(mockClasses, nil).
		Times(1)

	results, count, err := usecase.ListClasses(context.Background(), 10, 0, "id", "")

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)
	assert.Equal(t, "Economy", results[0].ClassName)
	assert.Equal(t, "Business", results[1].ClassName)
}

func TestClassUsecase_UpdateClass_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	updateRequest := &model.UpdateClassRequest{
		ID:         1,
		ClassName:  "Premium Economy",
		Type:       "Enhanced",
		ClassAlias: "PECO",
	}

	existingClass := &domain.Class{
		ID:         1,
		ClassName:  "Economy",
		Type:       "Standard",
		ClassAlias: "ECO",
	}

	mockClassRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingClass, nil).
		Times(1)

	mockClassRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, class *domain.Class) error {
			assert.Equal(t, uint(1), class.ID)
			assert.Equal(t, "Premium Economy", class.ClassName)
			assert.Equal(t, "Enhanced", class.Type)
			assert.Equal(t, "PECO", class.ClassAlias)
			return nil
		}).
		Times(1)

	err := usecase.UpdateClass(context.Background(), updateRequest)

	assert.NoError(t, err)
}

func TestClassUsecase_DeleteClass_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	db := setupClassTestDatabase()

	usecase := usecase.NewClassUsecase(db, mockClassRepo)

	existingClass := &domain.Class{
		ID:         1,
		ClassName:  "Economy",
		Type:       "Standard",
		ClassAlias: "ECO",
	}

	mockClassRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingClass, nil).
		Times(1)

	mockClassRepo.EXPECT().
		Delete(gomock.Any(), existingClass).
		Return(nil).
		Times(1)

	err := usecase.DeleteClass(context.Background(), 1)

	assert.NoError(t, err)
}
