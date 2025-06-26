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

type BookingUsecase struct {
	DB                *gorm.DB
	BookingRepository domain.BookingRepository
}

func NewBookingUsecase(
	db *gorm.DB,
	booking_repository domain.BookingRepository,
) *BookingUsecase {
	return &BookingUsecase{
		DB:                db,
		BookingRepository: booking_repository,
	}
}

func (b *BookingUsecase) CreateBooking(ctx context.Context, request *model.WriteBookingRequest) error {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	booking := &domain.Booking{
		OrderID:         request.OrderID,
		ScheduleID:      request.ScheduleID,
		CustomerName:    request.CustomerName,
		CustomerAge:     request.CustomerAge,
		CustomerGender:  request.CustomerGender,
		Email:           request.Email,
		PhoneNumber:     request.PhoneNumber,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		ReferenceNumber: request.ReferenceNumber,
	}
	if err := b.BookingRepository.Insert(tx, booking); err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (b *BookingUsecase) ListBookings(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadBookingResponse, int, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := b.BookingRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	bookings, err := b.BookingRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all bookings: %w", err)
	}

	responses := make([]*model.ReadBookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = mapper.BookingToResponse(booking)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (b *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	booking, err := b.BookingRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.BookingToResponse(booking), nil
}

func (b *BookingUsecase) GetBookingByOrderID(ctx context.Context, orderID string) (*model.ReadBookingResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	booking, err := b.BookingRepository.FindByOrderID(tx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.BookingToResponse(booking), nil
}

func (b *BookingUsecase) UpdateBooking(ctx context.Context, request *model.UpdateBookingRequest) error {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Fetch existing allocation
	booking, err := b.BookingRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return errs.ErrNotFound
	}

	booking.ScheduleID = request.ScheduleID
	booking.CustomerAge = request.CustomerAge
	booking.CustomerGender = request.CustomerGender
	booking.CustomerName = request.CustomerName
	booking.Email = request.Email
	booking.PhoneNumber = request.PhoneNumber
	booking.IDType = request.IDType
	booking.IDNumber = request.IDNumber
	booking.ReferenceNumber = request.ReferenceNumber
	booking.OrderID = request.OrderID

	if err := b.BookingRepository.Update(tx, booking); err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (b *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	booking, err := b.BookingRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return errs.ErrNotFound
	}

	if err := b.BookingRepository.Delete(tx, booking); err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
