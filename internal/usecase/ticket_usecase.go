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

type TicketUsecase struct {
	DB                 *gorm.DB
	TicketRepository   *repository.TicketRepository
	ScheduleRepository *repository.ScheduleRepository
	FareRepository     *repository.FareRepository
}

func NewTicketUsecase(
	db *gorm.DB,
	ticket_repository *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	fare_repository *repository.FareRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:                 db,
		TicketRepository:   ticket_repository,
		ScheduleRepository: schedule_repository,
		FareRepository:     fare_repository,
	}
}

func (t *TicketUsecase) CreateTicket(ctx context.Context, request *model.WriteTicketRequest) error {
	ticket := mapper.TicketMapper.FromWrite(request)

	if ticket.Status == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		return t.TicketRepository.Create(tx, ticket)
	})
}

func (t *TicketUsecase) GetAllTickets(ctx context.Context) ([]*model.ReadTicketResponse, error) {
	tickets := []*entity.Ticket{}

	err := tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		var err error
		tickets, err = t.TicketRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all tickets: %w", err)
	}

	return mapper.TicketMapper.ToModels(tickets), nil
}

func (t *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	ticket := new(entity.Ticket)

	err := tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		var err error
		ticket, err = t.TicketRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}

	if ticket == nil {
		return nil, errors.New("ticket not found")
	}

	return mapper.TicketMapper.ToModel(ticket), nil
}

func (t *TicketUsecase) UpdateTicket(ctx context.Context, id uint, request *model.UpdateTicketRequest) error {
	ticket := mapper.TicketMapper.FromUpdate(request)
	ticket.ID = id

	if ticket.ID == 0 {
		return fmt.Errorf("ticket ID cannot be zero")
	}

	if ticket.PassengerName == nil {
		return fmt.Errorf("passenger name cannot be empty")
	}

	return tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		return t.TicketRepository.Update(tx, ticket)
	})
}

func (t *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {

	return tx.Execute(ctx, t.DB, func(tx *gorm.DB) error {
		ticket, err := t.TicketRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if ticket == nil {
			return errors.New("ticket not found")
		}
		return t.TicketRepository.Delete(tx, ticket)
	})

}

// Execute handles the process of receiving and saving passenger data for claimed tickets,
// wrapped in a transaction.
func (t *TicketUsecase) FillData(ctx context.Context, request *model.FillPassengerDataRequest) (*model.FillPassengerDataResponse, error) {
	// Basic input validation (can remain outside the transaction)
	if len(request.PassengerData) == 0 {
		return nil, errors.New("invalid request: UserID and passenger data are required")
	}

	// Optional: Verify user exists (can remain outside the transaction if preferred,
	// or move inside if you need the user entity within the transaction)
	// _, err := uc.UserRepository.GetByID(ctx, request.UserID)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("user not found") }
	// 	return nil, fmt.Errorf("failed to verify user: %w", err)
	// }

	// Extract ticket IDs from the request (can remain outside the transaction)
	ticketIDs := make([]uint, len(request.PassengerData))
	passengerDataMap := make(map[uint]model.PassengerDataInput) // Map ticketID to submitted data
	for i, data := range request.PassengerData {
		ticketIDs[i] = data.TicketID
		passengerDataMap[data.TicketID] = data
	}

	// Variables to collect results from inside the transaction
	var updatedTicketIDs []uint
	var failedTickets []model.TicketUpdateFailure // Use the local struct type or model type consistently

	// --- Wrap the core logic in a transaction ---
	err := tx.Execute(ctx, t.DB, func(txDB *gorm.DB) error {
		// --- All repository calls within this function MUST use txDB ---

		// Retrieve the ticket entities from the database within the transaction
		// Use txDB for the repository call
		ticketsToUpdate, err := t.TicketRepository.FindManyByIDs(txDB, ticketIDs) // Use txDB
		if err != nil {
			// Log this error: Database read failure
			return fmt.Errorf("failed to retrieve tickets within transaction: %w", err) // Return error to trigger rollback
		}

		// Map retrieved tickets by ID for easy lookup
		retrievedTicketsMap := make(map[uint]*entity.Ticket)
		for _, ticket := range ticketsToUpdate {
			retrievedTicketsMap[ticket.ID] = ticket
		}

		// Prepare slice for bulk update if your repository supports UpdateMany
		ticketsForBulkUpdate := []*entity.Ticket{}
		now := time.Now() // Get time inside the transaction for consistency

		// --- Validate and Update Ticket Entities ---
		// Populate the slices declared outside the transaction func
		updatedTicketIDs = []uint{}
		failedTickets = []model.TicketUpdateFailure{}

		for _, reqData := range request.PassengerData {
			ticket, ok := retrievedTicketsMap[reqData.TicketID]
			if !ok {
				// Ticket ID from request was not found in the database
				failedTickets = append(failedTickets, model.TicketUpdateFailure{TicketID: reqData.TicketID, Reason: "Ticket not found"})
				continue
			}

			// Check Status
			if ticket.Status != "pending_data_entry" {
				failedTickets = append(failedTickets, model.TicketUpdateFailure{TicketID: reqData.TicketID, Reason: fmt.Sprintf("Ticket status is not pending data entry (%s)", ticket.Status)})
				continue
			}

			// Check Expiry Time
			if ticket.ExpiresAt.Before(now) { // Use 'now' captured inside the transaction
				failedTickets = append(failedTickets, model.TicketUpdateFailure{TicketID: reqData.TicketID, Reason: "Ticket has expired"})
				// Optional: Trigger cancellation of this ticket here or rely on background job
				continue
			}

			// Basic data validation (can be more extensive)
			if reqData.PassengerName == "" {
				failedTickets = append(failedTickets, model.TicketUpdateFailure{TicketID: reqData.TicketID, Reason: "Passenger name cannot be empty"})
				continue
			}
			// Add validation for PassportNumber, DateOfBirth, etc.

			// --- Update Fields if Valid ---
			// Update the entity with the submitted data
			ticket.PassengerName = &reqData.PassengerName // Assuming PassengerName in entity is *string
			ticket.IDType = &reqData.IDType               // Assuming PassportNumber in entity is *string
			ticket.IDNumber = &reqData.IDNumber           // Assuming PassportNumber in entity is *string
			ticket.SeatNumber = reqData.SeatNumber
			// ... update other passenger fields ...

			// Update status and timestamp
			ticket.Status = "pending_payment" // Move to next stage
			ticket.DataFilledAt = &now        // Set data filled timestamp (Assuming DataFilledAt is *time.Time)

			// Add to list for bulk update
			ticketsForBulkUpdate = append(ticketsForBulkUpdate, ticket)

			// Collect ID for successful updates list
			updatedTicketIDs = append(updatedTicketIDs, ticket.ID)
		}

		// --- Save Updated Tickets to Database ---
		// This step should be atomic for the batch.
		if len(ticketsForBulkUpdate) > 0 {
			// Assuming ITicketRepository has an UpdateMany or UpdateBulk method
			// Use txDB for the repository call
			err = t.TicketRepository.UpdateBulk(txDB, ticketsForBulkUpdate) // Use txDB
			if err != nil {
				// Log this error: Database write failure
				// If bulk update fails, return an error to trigger rollback of the entire transaction
				return fmt.Errorf("failed to save updated ticket data within transaction: %w", err)
			}
		}

		// If we reach here, all operations within the transaction function succeeded
		return nil // Return nil to trigger commit
	})
	// --- Transaction ends here (commit or rollback) ---

	// Handle any errors that occurred during the transaction
	if err != nil {
		// Check for specific errors returned from inside the transaction if needed
		// e.g., if using custom error types
		return nil, fmt.Errorf("failed to execute fill passenger data transaction: %w", err) // Wrap the error
	}

	// --- Return Response (outside the transaction) ---
	// The slices updated inside the transaction func are available here if commit succeeded
	return &model.FillPassengerDataResponse{
		UpdatedTicketIDs: updatedTicketIDs,
		FailedTickets:    failedTickets,
	}, nil
}
