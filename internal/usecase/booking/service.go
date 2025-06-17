package booking

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

type BookingUsecase struct {
	Tx                *tx.TxManager
	DB                *gorm.DB
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
}

func NewBookingUsecase(
	tx *tx.TxManager,
	db *gorm.DB,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
	session_repository *repository.SessionRepository,
) *BookingUsecase {
	return &BookingUsecase{
		Tx:                tx,
		DB:                db,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
		SessionRepository: session_repository,
	}
}

func (b *BookingUsecase) CreateBooking(ctx context.Context, request *model.WriteBookingRequest) error {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if request.CustomerName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	booking := &entity.Booking{
		OrderID:         request.OrderID,
		ScheduleID:      request.ScheduleID,
		CustomerName:    request.CustomerName,
		CustomerAge:     &request.CustomerAge,
		CustomerGender:  &request.CustomerGender,
		Email:           request.Email,
		PhoneNumber:     request.PhoneNumber,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		ReferenceNumber: request.ReferenceNumber,
		Status:          request.Status,
	}

	if err := b.BookingRepository.Create(tx, booking); err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (b *BookingUsecase) GetAllBookings(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadBookingResponse, int, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := b.BookingRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	bookings, err := b.BookingRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all bookings: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.BookingMapper.ToModels(bookings), int(total), nil
}

func (b *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	booking, err := b.BookingRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.BookingMapper.ToModel(booking), nil
}

func (b *BookingUsecase) GetBookingByOrderID(ctx context.Context, orderID string) (*model.ReadBookingResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	booking, err := b.BookingRepository.GetByOrderID(tx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.BookingMapper.ToModel(booking), nil
}

func (b *BookingUsecase) UpdateBooking(ctx context.Context, request *model.UpdateBookingRequest) error {
	tx := b.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	if request.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}

	// Fetch existing allocation
	booking, err := b.BookingRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return errors.New("booking not found")
	}

	booking.ScheduleID = request.ScheduleID
	booking.CustomerAge = &request.CustomerAge
	booking.CustomerGender = &request.CustomerGender
	booking.CustomerName = request.CustomerName
	booking.Email = request.Email
	booking.PhoneNumber = request.PhoneNumber
	booking.IDType = request.IDType
	booking.IDNumber = request.IDNumber
	booking.ReferenceNumber = request.ReferenceNumber
	booking.Status = request.Status
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
		} else {
			tx.Rollback()
		}
	}()

	booking, err := b.BookingRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return errors.New("booking not found")
	}

	if err := b.BookingRepository.Delete(tx, booking); err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (b *BookingUsecase) PaidConfirm(ctx context.Context, id uint) error {
	return b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return b.BookingRepository.PaidConfirm(tx, id)
	})
}
