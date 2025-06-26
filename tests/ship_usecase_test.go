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

func setupShipTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite: " + err.Error())
	}
	db.AutoMigrate(&domain.Ship{})
	return db
}

func TestShipUsecase_CreateShip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	req := &model.WriteShipRequest{
		ShipName:      "KM Budiono Siregar",
		ShipType:      "Ferry",
		ShipAlias:     "KMBS",
		Status:        "Active",
		YearOperation: "2020",
		ImageLink:     "https://example.com/ship.jpg",
		Description:   "Modern ferry for passenger transport",
	}

	mockShipRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, ship *domain.Ship) error {
			assert.Equal(t, "KM Budiono Siregar", ship.ShipName)
			assert.Equal(t, "Ferry", ship.ShipType)
			assert.Equal(t, "KMBS", ship.ShipAlias)
			assert.Equal(t, "Active", ship.Status)
			assert.Equal(t, "2020", ship.YearOperation)
			assert.Equal(t, "https://example.com/ship.jpg", ship.ImageLink)
			assert.Equal(t, "Modern ferry for passenger transport", ship.Description)
			return nil
		}).Times(1)

	err := usecase.CreateShip(context.Background(), req)

	assert.NoError(t, err)
}

func TestShipUsecase_CreateShip_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	req := &model.WriteShipRequest{
		ShipName:      "KM Budiono Siregar",
		ShipType:      "Ferry",
		ShipAlias:     "KMBS",
		Status:        "Active",
		YearOperation: "2020",
		ImageLink:     "https://example.com/ship.jpg",
		Description:   "Modern ferry for passenger transport",
	}

	mockShipRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(errors.New("database error")).Times(1)

	err := usecase.CreateShip(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create ship")
}

func TestShipUsecase_GetShipByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	expectedShip := &domain.Ship{
		ID:            1,
		ShipName:      "KM Budiono Siregar",
		ShipType:      "Ferry",
		ShipAlias:     "KMBS",
		Status:        "Active",
		YearOperation: "2020",
		ImageLink:     "https://example.com/ship.jpg",
		Description:   "Modern ferry for passenger transport",
	}

	mockShipRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(expectedShip, nil).
		Times(1)

	result, err := usecase.GetShipByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "KM Budiono Siregar", result.ShipName)
	assert.Equal(t, "Ferry", result.ShipType)
}

func TestShipUsecase_GetShipByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	mockShipRepo.EXPECT().
		FindByID(gomock.Any(), uint(999)).
		Return(nil, nil).
		Times(1)

	result, err := usecase.GetShipByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestShipUsecase_ListShips_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	mockShips := []*domain.Ship{
		{
			ID:            1,
			ShipName:      "KM Budiono Siregar",
			ShipType:      "Ferry",
			ShipAlias:     "KMBS",
			Status:        "Active",
			YearOperation: "2020",
		},
		{
			ID:            2,
			ShipName:      "KM Labobar",
			ShipType:      "Ferry",
			ShipAlias:     "KML",
			Status:        "Active",
			YearOperation: "2020",
		},
	}

	mockShipRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(2), nil).
		Times(1)

	mockShipRepo.EXPECT().
		FindAll(gomock.Any(), 10, 0, "id", "").
		Return(mockShips, nil).
		Times(1)

	results, count, err := usecase.ListShips(context.Background(), 10, 0, "id", "")

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)
	assert.Equal(t, "KM Budiono Siregar", results[0].ShipName)
	assert.Equal(t, "KM Labobar", results[1].ShipName)
}

func TestShipUsecase_UpdateShip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	updateRequest := &model.UpdateShipRequest{
		ID:            1,
		ShipName:      "KM Budiono Siregar Updated",
		ShipType:      "Fast Ferry",
		ShipAlias:     "KMBS-UPD",
		Status:        "Maintenance",
		YearOperation: "2021",
		ImageLink:     "https://example.com/ship-updated.jpg",
		Description:   "Updated modern ferry",
	}

	existingShip := &domain.Ship{
		ID:            1,
		ShipName:      "KM Budiono Siregar",
		ShipType:      "Ferry",
		ShipAlias:     "KMBS",
		Status:        "Active",
		YearOperation: "2020",
	}

	mockShipRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingShip, nil).
		Times(1)

	mockShipRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, ship *domain.Ship) error {
			assert.Equal(t, uint(1), ship.ID)
			assert.Equal(t, "KM Budiono Siregar Updated", ship.ShipName)
			assert.Equal(t, "Fast Ferry", ship.ShipType)
			assert.Equal(t, "KMBS-UPD", ship.ShipAlias)
			assert.Equal(t, "Maintenance", ship.Status)
			assert.Equal(t, "2021", ship.YearOperation)
			return nil
		}).
		Times(1)

	err := usecase.UpdateShip(context.Background(), updateRequest)

	assert.NoError(t, err)
}

func TestShipUsecase_DeleteShip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipRepo := mocks.NewMockShipRepository(ctrl)
	db := setupShipTestDatabase()

	usecase := usecase.NewShipUsecase(db, mockShipRepo)

	existingShip := &domain.Ship{
		ID:            1,
		ShipName:      "KM Budiono Siregar",
		ShipType:      "Ferry",
		ShipAlias:     "KMBS",
		Status:        "Active",
		YearOperation: "2020",
	}

	mockShipRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingShip, nil).
		Times(1)

	mockShipRepo.EXPECT().
		Delete(gomock.Any(), existingShip).
		Return(nil).
		Times(1)

	err := usecase.DeleteShip(context.Background(), 1)

	assert.NoError(t, err)
}
