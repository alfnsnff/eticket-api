package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"eticket-api/pkg/gotann"
	"fmt"
)

type BookingUsecase struct {
	Transactor        *transact.Transactor
	BookingRepository domain.BookingRepository
}

func NewBookingUsecase(
	transactor *transact.Transactor,
	booking_repository domain.BookingRepository,
) *BookingUsecase {
	return &BookingUsecase{
		Transactor:        transactor,
		BookingRepository: booking_repository,
	}
}

func (uc *BookingUsecase) CreateBooking(ctx context.Context, request *model.WriteBookingRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking := &domain.Booking{
			OrderID:         request.OrderID,
			ScheduleID:      request.ScheduleID,
			IDType:          request.IDType,
			IDNumber:        request.IDNumber,
			CustomerName:    request.CustomerName,
			CustomerAge:     request.CustomerAge,
			CustomerGender:  request.CustomerGender,
			PhoneNumber:     request.PhoneNumber,
			Email:           request.Email,
			ReferenceNumber: request.ReferenceNumber,
		}
		if err := uc.BookingRepository.Insert(ctx, tx, booking); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create booking: %w", err)
		}
		return nil
	})
}

func (uc *BookingUsecase) ListBookings(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadBookingResponse, int, error) {
	var err error
	var total int64
	var bookings []*domain.Booking
	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.BookingRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count bookings: %w", err)
		}

		bookings, err = uc.BookingRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all bookings: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}

	responses := make([]*model.ReadBookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = mapper.BookingToResponse(booking)
	}

	return responses, int(total), nil
}

func (uc *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	var err error
	var booking *domain.Booking
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound

		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get booking by id: %w", err)
	}

	return mapper.BookingToResponse(booking), nil
}

func (uc *BookingUsecase) GetBookingByOrderID(ctx context.Context, orderID string) (*model.ReadBookingResponse, error) {
	var err error
	var booking *domain.Booking
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err = uc.BookingRepository.FindByOrderID(ctx, tx, orderID)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get booking by order ID: %w", err)
	}

	return mapper.BookingToResponse(booking), nil
}

func (uc *BookingUsecase) UpdateBooking(ctx context.Context, request *model.UpdateBookingRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByID(ctx, tx, request.ID)
		if err != nil {
			return fmt.Errorf("failed to find booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}

		booking.OrderID = request.OrderID
		booking.ScheduleID = request.ScheduleID
		booking.IDType = request.IDType
		booking.IDNumber = request.IDNumber
		booking.CustomerName = request.CustomerName
		booking.CustomerAge = request.CustomerAge
		booking.CustomerGender = request.CustomerGender
		booking.PhoneNumber = request.PhoneNumber
		booking.Email = request.Email
		booking.ReferenceNumber = request.ReferenceNumber

		if err := uc.BookingRepository.Update(ctx, tx, booking); err != nil {
			return fmt.Errorf("failed to update booking: %w", err)
		}
		return nil
	})
}

func (uc *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}

		if err := uc.BookingRepository.Delete(ctx, tx, booking); err != nil {
			return fmt.Errorf("failed to delete booking: %w", err)
		}
		return nil
	})
}
