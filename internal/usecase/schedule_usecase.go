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

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	DB                   *gorm.DB
	AllocationRepository *repository.AllocationRepository
	ClassRepository      *repository.ClassRepository
	FareRepository       *repository.FareRepository
	ManifestRepository   *repository.ManifestRepository
	RouteRepository      *repository.RouteRepository
	ShipRepository       *repository.ShipRepository
	ScheduleRepository   *repository.ScheduleRepository
	TicketRepository     *repository.TicketRepository
}

func NewScheduleUsecase(
	db *gorm.DB,
	allocation_repository *repository.AllocationRepository,
	class_repository *repository.ClassRepository,
	fare_repository *repository.FareRepository,
	manifest_repository *repository.ManifestRepository,
	route_repository *repository.RouteRepository,
	ship_repository *repository.ShipRepository,
	schedule_repository *repository.ScheduleRepository,
	ticket_repository *repository.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		DB:                   db,
		AllocationRepository: allocation_repository,
		ClassRepository:      class_repository,
		FareRepository:       fare_repository,
		ManifestRepository:   manifest_repository,
		RouteRepository:      route_repository,
		ShipRepository:       ship_repository,
		ScheduleRepository:   schedule_repository,
		TicketRepository:     ticket_repository,
	}
}

func (sc *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.WriteScheduleRequest) error {
	schedule := mapper.ScheduleMapper.FromWrite(request)

	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		return sc.ScheduleRepository.Create(tx, schedule)
	})
}

func (sc *ScheduleUsecase) GetAllSchedules(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	schedules := []*entity.Schedule{}

	err := tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		var err error
		schedules, err = sc.ScheduleRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	return mapper.ScheduleMapper.ToModels(schedules), nil
}

func (sc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	schedule := new(entity.Schedule)

	err := tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		var err error
		schedule, err = sc.ScheduleRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}

	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	return mapper.ScheduleMapper.ToModel(schedule), nil
}

func (sc *ScheduleUsecase) UpdateSchedule(ctx context.Context, id uint, request *model.UpdateScheduleRequest) error {
	schedule := mapper.ScheduleMapper.FromUpdate(request)
	schedule.ID = id

	if schedule.ID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}

	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		return sc.ScheduleRepository.Update(tx, schedule)
	})
}

func (sc *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {

	return tx.Execute(ctx, sc.DB, func(tx *gorm.DB) error {
		schedule, err := sc.ScheduleRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}
		return sc.ScheduleRepository.Delete(tx, schedule)
	})

}

// // SearchSchedule searches a schedule by departure, arrival, and date
// func (s *ScheduleUsecase) SearchSchedule(ctx context.Context, request *model.ScheduleSearchRequest) (*entity.Schedule, error) {
// 	var schedule *entity.Schedule

// 	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
// 		schedule, err = s.ScheduleRepository.Search(txDB, route.ID, req.Date, req.ShipID)
// 		return err
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search schedule: %w", err)
// 	}

// 	return schedule, nil
// }

// Refactored to use tx.Execute helper for demonstrating transaction usage even on reads.
func (sc *ScheduleUsecase) GetScheduleDetailsWithAvailability(ctx context.Context, scheduleID uint) (*model.ReadScheduleDetailsWithAvailabilityResponse, error) {

	var schedule *entity.Schedule
	var scheduleCapacities []*entity.Allocation
	classesAvailability := make([]model.ScheduleClassAvailability, 0) // Initialize outside to be populated inside

	// --- Wrap all database interactions inside tx.Execute ---
	err := tx.Execute(ctx, sc.DB, func(txDB *gorm.DB) error {
		// --- All repository calls within this function MUST use txDB ---

		// 1. Get the basic Schedule details
		// Use txDB instead of sc.ScheduleRepository.DB
		var err error                                                   // Declare err locally within the transaction function
		schedule, err = sc.ScheduleRepository.GetByID(txDB, scheduleID) // Use txDB
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Return a specific error that tx.Execute will wrap
				return errors.New("schedule not found") // Return error to trigger rollback
			}
			return fmt.Errorf("failed to get schedule by ID within transaction: %w", err) // Return error to trigger rollback
		}

		// // Get related Ship and Route details for the response
		// ship, err = sc.ShipRepository.GetByID(txDB, schedule.ShipID) // Use txDB
		// if err != nil {
		// 	// Log this error: Data inconsistency
		// 	return fmt.Errorf("failed to get ship details within transaction: %w", err) // Return error
		// }

		// 2. Get all ScheduleCapacity entries for this schedule
		// Use txDB instead of sc.ScheduleCapacityRepository.DB
		scheduleCapacities, err = sc.AllocationRepository.FindByScheduleID(txDB, scheduleID) // Use txDB
		if err != nil {
			return fmt.Errorf("failed to get schedule capacities within transaction: %w", err) // Return error
		}

		// Prepare the slice for class availability details
		// classesAvailability = make([]model.ScheduleClassAvailability, 0, len(scheduleCapacities)) // Can be done outside or here

		// 3. For each ScheduleCapacity entry (each class on this schedule):
		for _, scap := range scheduleCapacities {
			// 3a. Count occupied tickets for this specific schedule and class
			// Use txDB instead of sc.TicketRepository.DB
			// No FOR UPDATE needed for this read count
			occupiedCount, err := sc.TicketRepository.CountByScheduleClassAndStatuses(txDB, scheduleID, scap.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"}) // Use txDB
			if err != nil {
				// Log this error: Database issue counting tickets
				return fmt.Errorf("failed to count occupied tickets within transaction: %w", err) // Return error
			}

			// 3b. Calculate available capacity
			availableCapacity := scap.Quota - int(occupiedCount)

			// 3c. Get Class name
			// Use txDB instead of sc.ClassRepository.DB
			class, err := sc.ClassRepository.GetByID(txDB, scap.ClassID) // Use txDB
			if err != nil {
				// Log this error: Data inconsistency
				return fmt.Errorf("failed to get class details within transaction: %w", err) // Return error
			}
			if class == nil {
				// Log this error: Data inconsistency
				return fmt.Errorf("class not found for ID %d within transaction", scap.ClassID) // Return error
			}

			// 3d. Look up price using Manifest and Fare
			// Use txDB instead of sc.ManifestRepository.DB / sc.FareRepository.DB
			manifest, err := sc.ManifestRepository.GetByShipAndClass(txDB, schedule.ShipID, scap.ClassID) // Use txDB
			if err != nil {
				// Log this error: Configuration error?
				return fmt.Errorf("failed to get manifest within transaction: %w", err) // Return error
			}
			if manifest == nil {
				// Log this error: Configuration error?
				return fmt.Errorf("manifest not found for ship %d, class %d within transaction", schedule.ShipID, scap.ClassID) // Return error
			}

			fare, err := sc.FareRepository.GetByManifestAndRoute(txDB, manifest.ID, schedule.RouteID) // Use txDB
			if err != nil {
				// Log this error: Configuration error?
				return fmt.Errorf("failed to get fare within transaction: %w", err) // Return error
			}
			if fare == nil {
				// Log this error: Configuration error?
				return fmt.Errorf("fare not found for manifest %d, route %d within transaction", manifest.ID, schedule.RouteID) // Return error
			}

			// 3e. Add to the results slice (can be done inside the transaction func)
			classesAvailability = append(classesAvailability, model.ScheduleClassAvailability{
				ClassID:           scap.ClassID,
				ClassName:         class.Name,
				TotalCapacity:     scap.Quota,
				AvailableCapacity: availableCapacity,
				Price:             fare.Price,
				Currency:          "IDR", // Or get from config/Fare entity
			})
		}

		// Return nil if all operations within the transaction function succeeded
		return nil
	})
	// --- Transaction ends here (commit or rollback) ---

	// Handle any errors that occurred during the transaction
	if err != nil {
		// Check for the specific "schedule not found" error returned from inside the transaction
		if errors.New("schedule not found").Error() == err.Error() { // Compare error strings or use custom error types
			return nil, errors.New("schedule not found") // Re-return the user-friendly error
		}
		return nil, fmt.Errorf("failed to get schedule details with availability transaction: %w", err) // Wrap other errors
	}

	// 4. Build the final response model (outside the transaction)
	// The variables populated inside the transaction func (schedule, ship, route, classesAvailability)
	// are available here if the transaction committed successfully.
	response := &model.ReadScheduleDetailsWithAvailabilityResponse{
		ScheduleID:          schedule.ID,
		RouteID:             schedule.RouteID,
		ClassesAvailability: classesAvailability,
	}

	return response, nil // Return the response and nil error
}
