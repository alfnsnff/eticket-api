package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PrebookTicketsUsecase handles the logic for claiming ticket slots
type PrebookTicketsUsecase struct {
	DB                   *gorm.DB // Database instance for transaction management
	ScheduleRepository   *repository.ScheduleRepository
	AllocationRepository *repository.AllocationRepository // Your AllocationRepository implements this
	TicketRepository     *repository.TicketRepository
	ManifestRepository   *repository.ManifestRepository
	FareRepository       *repository.FareRepository
}

// PrebookTicketsUsecase creates a new PrebookTicketsUsecase
func NewPrebookTicketsUsecase(
	db *gorm.DB,
	schedule_repository *repository.ScheduleRepository,
	allocation_repository *repository.AllocationRepository,
	ticket_repository *repository.TicketRepository,
	manifest_repository *repository.ManifestRepository,
	fare_repository *repository.FareRepository,
) *PrebookTicketsUsecase {
	return &PrebookTicketsUsecase{
		DB:                   db,
		ScheduleRepository:   schedule_repository,
		AllocationRepository: allocation_repository,
		TicketRepository:     ticket_repository,
		ManifestRepository:   manifest_repository,
		FareRepository:       fare_repository,
	}
}

// Execute handles the process of claiming ticket slots atomically
func (uc *PrebookTicketsUsecase) Execute(ctx context.Context, request *model.ClaimTicketsRequest) (*model.ClaimTicketsResponse, error) {
	// Input validation
	if request.ScheduleID == 0 || len(request.Items) == 0 {
		return nil, errors.New("invalid claim request")
	}

	var ClaimedTicketIDs []uint
	var expiryTime time.Time

	// Use the transaction helper
	err := tx.Execute(ctx, uc.DB, func(tx *gorm.DB) error {
		// This function runs within the transaction
		// Use the repositories with the transaction instance (tx)

		// Map to store availability checks results
		availabilityChecks := make(map[uint]int64) // classID -> availableCount

		// --- STEP 1: Lock Capacity Rows and Check Availability for ALL Requested Classes ---
		for _, item := range request.Items {
			if item.Quantity == 0 {
				continue // Skip items with 0 quantity
			}

			// Get the ScheduleCapacity entity WITH a FOR UPDATE lock
			// Use the repository method we defined aerlier
			scheduleCapacity, err := uc.AllocationRepository.LockByScheduleAndClass(tx, request.ScheduleID, item.ClassID)
			if err != nil {
				return fmt.Errorf("failed to get schedule capacity for class %d: %w", item.ClassID, err)
			}
			if scheduleCapacity == nil {
				// This schedule/class combination doesn't exist in ScheduleCapacity - configuration error?
				return fmt.Errorf("allocation ot found for schedule %d, class %d", request.ScheduleID, item.ClassID)
			}

			// Count currently occupied tickets for this schedule/class
			occupiedCount, err := uc.TicketRepository.CountByScheduleClassAndStatuses(tx, request.ScheduleID, item.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"})
			if err != nil {
				return fmt.Errorf("failed to count occupied tickets for class %d: %w", item.ClassID, err)
			}

			// Calculate available count and store it
			availableCount := int64(scheduleCapacity.Quota) - occupiedCount
			availabilityChecks[item.ClassID] = availableCount

			// Check if the requested quantity exceeds available
			if availableCount < int64(item.Quantity) {
				// IMPORTANT: If any class fails, return an error to ROLLBACK the entire transaction
				return fmt.Errorf("not enough available slots for class %d. Available: %d, Requested: %d",
					item.ClassID, availableCount, item.Quantity)
			}
		}

		// --- STEP 2: If ALL Checks Pass, Create Ticket Entities and Insert ---
		// This part only runs if the loop above completed without returning an error
		now := time.Now()
		expiryTime = now.Add(15 * time.Minute) // Define your timeout duration

		ticketsToCreate := []*entity.Ticket{}
		claimedTicketIDs := []uint{} // Collect IDs for the response

		// Get Schedule details to fetch ShipID and RouteID for Fare lookup
		schedule, err := uc.ScheduleRepository.GetByID(tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("failed to get schedule details: %w", err) // Should not happen if capacity existed
		}

		for _, item := range request.Items {
			if item.Quantity == 0 {
				continue
			}

			// Get Manifest ID for price lookup
			manifest, err := uc.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, item.ClassID)
			if err != nil {
				return fmt.Errorf("failed to get manifest for ship %d, class %d: %w", schedule.ShipID, item.ClassID, err)
			}
			if manifest == nil {
				// Configuration error: Manifest not defined for this ship/class?
				return fmt.Errorf("manifest not found for ship %d, class %d", schedule.ShipID, item.ClassID)
			}

			// Get Price for this Manifest and Schedule's Route
			fare, err := uc.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID) // Assuming Manifest entity has an ID field
			if err != nil {
				return fmt.Errorf("failed to get fare for manifest %d, route %d: %w", manifest.ID, schedule.RouteID, err)
			}
			if fare == nil {
				// Configuration error: Fare not defined for this combination?
				return fmt.Errorf("fare not found for manifest %d, route %d", manifest.ID, schedule.RouteID)
			}

			for i := 0; i < int(item.Quantity); i++ {
				ticket := &entity.Ticket{
					// ID will be generated by the database
					ScheduleID: request.ScheduleID,
					ClassID:    item.ClassID,
					Status:     "pending_data_entry", // Initial status
					Price:      fare.Price,           // Store the price at the time of claim
					ClaimedAt:  now,
					ExpiresAt:  expiryTime,
					// BookingID is NULL initially
					// PassengerName, SeatNumber, DataFilledAt, BookingTimestamp are NULL initially
				}
				ticketsToCreate = append(ticketsToCreate, ticket)
			}
		}

		// Insert all the new ticket records in one go
		// Assuming ITicketRepository.CreateMany takes []*entity.Ticket
		err = uc.TicketRepository.CreateBulk(tx, ticketsToCreate)
		if err != nil {
			return fmt.Errorf("failed to create tickets: %w", err) // This will trigger rollback
		}

		// Collect the generated IDs from the created tickets for the response
		// GORM should populate the IDs after CreateMany
		for _, ticket := range ticketsToCreate {
			claimedTicketIDs = append(claimedTicketIDs, ticket.ID)
		}

		ClaimedTicketIDs = claimedTicketIDs

		// If we reach here, the transaction will be committed by tx.Execute
		return nil // No error means commit
	})

	// This runs after the transaction is complete (commit or rollback)
	if err != nil {
		// Handle specific error types from the availability checks if needed
		// e.g., check if errors.Is(err, errors.New("not enough available slots..."))
		return nil, fmt.Errorf("failed to claim tickets: %w", err)
	}

	// Return the response model
	return &model.ClaimTicketsResponse{
		ClaimedTicketIDs: ClaimedTicketIDs,
		ExpiresAt:        expiryTime,
	}, nil
}
