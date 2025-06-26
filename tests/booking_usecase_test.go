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

func setupBookingTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite: " + err.Error())
	}
	db.AutoMigrate(&domain.Booking{})
	return db
}

func TestBookingUsecase_CreateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	orderID := "ORDER123"
	referenceNumber := "REF123"
	req := &model.WriteBookingRequest{
		OrderID:         &orderID,
		ScheduleID:      1,
		CustomerName:    "John Doe",
		CustomerAge:     30,
		CustomerGender:  "Male",
		Email:           "john@example.com",
		PhoneNumber:     "081234567890",
		IDType:          "KTP",
		IDNumber:        "1234567890123456",
		ReferenceNumber: &referenceNumber,
	}
	mockBookingRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, booking *domain.Booking) error {
			assert.Equal(t, "ORDER123", *booking.OrderID)
			assert.Equal(t, uint(1), booking.ScheduleID)
			assert.Equal(t, "John Doe", booking.CustomerName)
			assert.Equal(t, "john@example.com", booking.Email)
			return nil
		}).Times(1)

	err := usecase.CreateBooking(context.Background(), req)

	assert.NoError(t, err)
}

func TestBookingUsecase_CreateBooking_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	orderID := "ORDER123"
	referenceNumber := "REF123"
	req := &model.WriteBookingRequest{
		OrderID:         &orderID,
		ScheduleID:      1,
		CustomerName:    "John Doe",
		CustomerAge:     30,
		CustomerGender:  "Male",
		Email:           "john@example.com",
		PhoneNumber:     "081234567890",
		IDType:          "KTP",
		IDNumber:        "1234567890123456",
		ReferenceNumber: &referenceNumber,
	}

	mockBookingRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(errors.New("database error")).Times(1)

	err := usecase.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create booking")
}

func TestBookingUsecase_GetBookingByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	orderID := "ORDER123"
	referenceNumber := "REF123"
	expectedBooking := &domain.Booking{
		ID:              1,
		OrderID:         &orderID,
		ScheduleID:      1,
		CustomerName:    "John Doe",
		CustomerAge:     30,
		CustomerGender:  "Male",
		Email:           "john@example.com",
		PhoneNumber:     "081234567890",
		IDType:          "KTP",
		IDNumber:        "1234567890123456",
		ReferenceNumber: &referenceNumber,
	}
	mockBookingRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(expectedBooking, nil).
		Times(1)

	result, err := usecase.GetBookingByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "ORDER123", *result.OrderID)
	assert.Equal(t, "John Doe", result.CustomerName)
}

func TestBookingUsecase_GetBookingByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	mockBookingRepo.EXPECT().
		FindByID(gomock.Any(), uint(999)).
		Return(nil, nil).
		Times(1)

	result, err := usecase.GetBookingByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestBookingUsecase_ListBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	orderID1 := "ORDER123"
	orderID2 := "ORDER123"
	mockBookings := []*domain.Booking{
		{
			ID:           1,
			OrderID:      &orderID1,
			CustomerName: "John Doe",
		},
		{
			ID:           2,
			OrderID:      &orderID2,
			CustomerName: "Jane Doe",
		},
	}

	mockBookingRepo.EXPECT().
		Count(gomock.Any()).
		Return(int64(2), nil).
		Times(1)

	mockBookingRepo.EXPECT().
		FindAll(gomock.Any(), 10, 0, "id", "").
		Return(mockBookings, nil).
		Times(1)

	results, count, err := usecase.ListBookings(context.Background(), 10, 0, "id", "")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)
	assert.Equal(t, "ORDER123", *results[0].OrderID)
}

func TestBookingUsecase_UpdateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	orderID1 := "ORDER122"
	orderID2 := "ORDER124"
	referenceNumber := "REF123"
	updateRequest := &model.UpdateBookingRequest{
		ID:              1,
		ScheduleID:      2,
		CustomerName:    "John Updated",
		CustomerAge:     31,
		CustomerGender:  "Male",
		Email:           "john.updated@example.com",
		PhoneNumber:     "081234567891",
		IDType:          "KTP",
		IDNumber:        "1234567890123457",
		ReferenceNumber: &referenceNumber,
		OrderID:         &orderID2,
	}

	existingBooking := &domain.Booking{
		ID:           1,
		OrderID:      &orderID1,
		CustomerName: "John Doe",
	}

	mockBookingRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingBooking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, booking *domain.Booking) error {
			assert.Equal(t, uint(1), booking.ID)
			assert.Equal(t, uint(2), booking.ScheduleID)
			assert.Equal(t, "John Updated", booking.CustomerName)
			assert.Equal(t, "john.updated@example.com", booking.Email)
			return nil
		}).
		Times(1)

	err := usecase.UpdateBooking(context.Background(), updateRequest)

	assert.NoError(t, err)
}

func TestBookingUsecase_DeleteBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	db := setupBookingTestDatabase()

	usecase := usecase.NewBookingUsecase(db, mockBookingRepo)

	orderID := "ORDER123"
	existingBooking := &domain.Booking{
		ID:           1,
		OrderID:      &orderID,
		CustomerName: "John Doe",
	}

	mockBookingRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(existingBooking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		Delete(gomock.Any(), existingBooking).
		Return(nil).
		Times(1)

	err := usecase.DeleteBooking(context.Background(), 1)

	assert.NoError(t, err)
}
