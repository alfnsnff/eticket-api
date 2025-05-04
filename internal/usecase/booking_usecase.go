package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookingUsecase struct {
	DB                *gorm.DB
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
}

func NewBookingUsecase(db *gorm.DB,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
) *BookingUsecase {
	return &BookingUsecase{
		DB:                db,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
	}
}

func (b *BookingUsecase) CreateBooking(ctx context.Context, request *model.WriteBookingRequest) error {
	booking := mapper.BookingMapper.FromWrite(request)

	if booking.CustomerName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		return b.BookingRepository.Create(tx, booking)
	})
}

func (b *BookingUsecase) GetAllBookings(ctx context.Context) ([]*model.ReadBookingResponse, error) {
	bookings := []*entity.Booking{}

	err := tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		var err error
		bookings, err = b.BookingRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.BookingMapper.ToModels(bookings), nil
}

func (b *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	booking := new(entity.Booking)

	err := tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		var err error
		booking, err = b.BookingRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, errors.New("booking not found")
	}

	return mapper.BookingMapper.ToModel(booking), nil
}

func (b *BookingUsecase) UpdateBooking(ctx context.Context, id uint, request *model.UpdateBookingRequest) error {
	booking := mapper.BookingMapper.FromUpdate(request)
	booking.ID = id

	if booking.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}
	if booking.CustomerName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		return b.BookingRepository.Update(tx, booking)
	})
}

// Deletebooking deletes a booking by its ID
func (b *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {

	return tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		booking, err := b.BookingRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if booking == nil {
			return errors.New("route not found")
		}
		return b.BookingRepository.Delete(tx, booking)
	})

}

// Execute handles the process of confirming payment and finalizing the booking.
// This is a critical atomic transaction.
func (b *BookingUsecase) ConfirmBook(ctx context.Context, request *model.ConfirmPaymentRequest) (*model.ConfirmPaymentResponse, error) {
	// Basic input validation
	if len(request.TicketIDs) == 0 ||
		request.Name == "" || request.IDType == "" || request.IDNumber == "" ||
		request.PhoneNumber == "" || request.Email == "" || request.BirthDate.IsZero() {
		return nil, errors.New("invalid request: PaymentIntentID, TicketIDs, and all booker details are required")
	}

	// User verification step is removed as users are not logged in

	var confirmedBooking *entity.Booking
	var confirmedTicketIDs []uint

	// --- Wrap the core logic in a transaction ---
	err := tx.Execute(ctx, b.DB, func(txDB *gorm.DB) error {
		// --- All repository calls within this function MUST use txDB ---

		// Step 1: Verify Payment Status (Often involves an external call)
		// This part depends on your payment gateway integration.
		// You might call an external service here or check a local payment record status.
		// For this example, we'll assume a successful payment check happened already
		// or is handled by the PaymentIntentID verification.
		// If payment verification fails, return an error to trigger rollback.
		// Example: Check payment record status in your DB if you store them.
		// paymentRecord, err := uc.PaymentRepository.GetByIntentID(txDB, request.PaymentIntentID) // Use txDB
		// if err != nil || paymentRecord.Status != "succeeded" {
		//    return errors.New("payment not verified or failed")
		// }

		// Step 2: Retrieve and Verify Tickets
		// Retrieve the tickets that are supposed to be confirmed by this payment.
		// Use txDB for the repository call.
		ticketsToConfirm, err := b.TicketRepository.FindManyByIDs(txDB, request.TicketIDs) // Use txDB
		if err != nil {
			// Log this error: Database read failure
			return fmt.Errorf("failed to retrieve tickets for confirmation within transaction: %w", err) // Return error
		}

		if len(ticketsToConfirm) != len(request.TicketIDs) {
			// Some requested tickets were not found. This might indicate a problem.
			// You might want to return an error or handle this case specifically.
			return errors.New("one or more tickets for confirmation not found")
		}

		var totalAmount float32 = 0
		ticketsForBulkUpdate := []*entity.Ticket{}
		now := time.Now() // Get time inside the transaction

		// Verify status and expiry for each ticket
		// Also, get the ScheduleID from one of the tickets (assuming all are for the same schedule)
		var scheduleID uint = 0
		if len(ticketsToConfirm) > 0 {
			// Assuming ScheduleID on Ticket entity is NOT nullable (uint)
			scheduleID = ticketsToConfirm[0].ScheduleID // Assuming all tickets are for the same schedule
		} else {
			// Should not happen if len(request.TicketIDs) > 0 and FindManyByIDs returned a slice of the same length
			return errors.New("no tickets found to confirm")
		}

		for _, ticket := range ticketsToConfirm {
			// Check Status - Must be in a state ready for confirmation (e.g., pending_payment)
			if ticket.Status != "pending_payment" {
				// Log this: Ticket in wrong state for confirmation
				return fmt.Errorf("ticket %d in wrong status for confirmation: %s", ticket.ID, ticket.Status) // Return error
			}

			// Check Expiry Time - Ensure it hasn't expired just before confirmation
			if ticket.ExpiresAt.Before(now) { // Use 'now' captured inside the transaction
				// Log this: Ticket expired just before confirmation
				return fmt.Errorf("ticket %d expired just before confirmation", ticket.ID) // Return error
				// The background job should ideally cancel this, but this is a final check
			}

			// Accumulate total amount
			totalAmount += ticket.Price

			// Prepare for update
			ticket.Status = "confirmed"
			ticket.BookingTimestamp = &now // Set confirmation timestamp (Assuming BookingTimestamp on Ticket is *time.Time)
			// BookingID will be set after the Booking entity is created
			ticketsForBulkUpdate = append(ticketsForBulkUpdate, ticket)
		}

		// Step 3: Create the Booking Record using the provided entity structure
		// Use txDB for the repository call.
		newBooking := &entity.Booking{
			ScheduleID: scheduleID, // Populate ScheduleID from tickets

			// Populate Booker Information from the request
			CustomerName: request.Name,
			IDType:       request.IDType,
			IDNumber:     request.IDNumber,
			PhoneNumber:  request.PhoneNumber,
			Email:        request.Email,
			BirthDate:    request.BirthDate,

			// Populate Booking Transaction Details
			BookingTimestamp: now, // Use the same timestamp as ticket confirmation
			TotalAmount:      totalAmount,
			Status:           "completed", // Or "paid", "confirmed"

			// Add other booking details like PaymentIntentID if needed
			// PaymentIntentID: request.PaymentIntentID,
		}
		err = b.BookingRepository.Create(txDB, newBooking) // Use txDB
		if err != nil {
			// Log this error: Database write failure
			return fmt.Errorf("failed to create booking record within transaction: %w", err) // Return error
		}

		// Store the created booking entity to access its ID later
		confirmedBooking = newBooking

		// Step 4: Update Tickets with Booking ID and Status
		// Now that we have the BookingID, update the tickets
		confirmedTicketIDs = make([]uint, len(ticketsForBulkUpdate)) // Populate slice declared outside
		for i, ticket := range ticketsForBulkUpdate {
			// Assuming BookingID on Ticket entity is *uint
			ticket.BookingID = &confirmedBooking.ID // Set the FK to the new Booking ID
			confirmedTicketIDs[i] = ticket.ID       // Collect IDs for response
		}

		// Save the updated tickets (now with BookingID and status='confirmed')
		// Assuming ITicketRepository has an UpdateMany method
		err = b.TicketRepository.UpdateBulk(txDB, ticketsForBulkUpdate) // Use txDB
		if err != nil {
			// Log this error: Database write failure
			return fmt.Errorf("failed to update tickets with booking ID within transaction: %w", err) // Return error
		}

		// If we reach here, all operations within the transaction function succeeded
		return nil // Return nil to trigger commit
	})
	// --- Transaction ends here (commit or rollback) ---

	// Handle any errors that occurred during the transaction
	if err != nil {
		// Check for specific errors returned from inside the transaction if needed
		// e.g., "ticket in wrong status", "ticket expired", "payment not verified"
		return nil, fmt.Errorf("failed to execute payment confirmation transaction: %w", err) // Wrap the error
	}

	// --- Return Response (outside the transaction) ---
	// The variables populated inside the transaction func are available here if commit succeeded
	return &model.ConfirmPaymentResponse{
		BookingID:          confirmedBooking.ID,
		BookingStatus:      confirmedBooking.Status,
		ConfirmedTicketIDs: confirmedTicketIDs,
	}, nil
}
