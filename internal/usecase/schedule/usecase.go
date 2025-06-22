package schedule

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	DB                   *gorm.DB // Assuming you have a DB field for the transaction manager
	AllocationRepository AllocationRepository
	ClassRepository      ClassRepository
	FareRepository       FareRepository
	ManifestRepository   ManifestRepository
	ShipRepository       ShipRepository
	ScheduleRepository   ScheduleRepository
	TicketRepository     TicketRepository
}

func NewScheduleUsecase(
	db *gorm.DB,
	allocation_repository AllocationRepository,
	class_repository ClassRepository,
	fare_repository FareRepository,
	manifest_repository ManifestRepository,
	ship_repository ShipRepository,
	schedule_repository ScheduleRepository,
	ticket_repository TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		DB:                   db,
		AllocationRepository: allocation_repository,
		ClassRepository:      class_repository,
		FareRepository:       fare_repository,
		ManifestRepository:   manifest_repository,
		ShipRepository:       ship_repository,
		ScheduleRepository:   schedule_repository,
		TicketRepository:     ticket_repository,
	}
}

func (sc *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.WriteScheduleRequest) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule := &domain.Schedule{
		RouteID:           request.RouteID,
		ShipID:            request.ShipID,
		DepartureDatetime: &request.DepartureDatetime,
		ArrivalDatetime:   &request.ArrivalDatetime,
		Status:            &request.Status,
	}

	if err := sc.ScheduleRepository.Create(tx, schedule); err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sc *ScheduleUsecase) GetAllSchedules(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadScheduleResponse, int, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := sc.AllocationRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	schedules, err := sc.ScheduleRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all schedules: %w", err)
	}

	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedules := range schedules {
		responses[i] = ScheduleToResponse(schedules)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (sc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule, err := sc.ScheduleRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ScheduleToResponse(schedule), nil
}

func (sc *ScheduleUsecase) UpdateSchedule(ctx context.Context, request *model.UpdateScheduleRequest) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	schedule, err := sc.ScheduleRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find schedule: %w", err)
	}
	if schedule == nil {
		return errors.New("schedule not found")
	}

	schedule.RouteID = request.RouteID
	schedule.ShipID = request.ShipID
	schedule.DepartureDatetime = &request.DepartureDatetime
	schedule.ArrivalDatetime = &request.ArrivalDatetime
	schedule.Status = &request.Status

	if err := sc.ScheduleRepository.Update(tx, schedule); err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sc *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	schedule, err := sc.ScheduleRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule not found")
	}

	if err := sc.ScheduleRepository.Delete(tx, schedule); err != nil {
		return fmt.Errorf("failed to delete allocation: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sc *ScheduleUsecase) GetAllScheduled(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	schedules, err := sc.ScheduleRepository.GetActiveSchedule(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all active schedules: %w", err)
	}

	responses := make([]*model.ReadScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = ScheduleToResponse(schedule)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, nil
}
func (sc *ScheduleUsecase) GetScheduleAvailability(ctx context.Context, scheduleID uint) (*model.ReadScheduleDetailsResponse, error) {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback() // read-only logic; rollback after use
		}
	}()

	schedule, err := sc.ScheduleRepository.GetByID(tx, scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	scheduleCapacities, err := sc.AllocationRepository.FindByScheduleID(tx, schedule.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule capacities: %w", err)
	}

	classAvailabilities := make([]model.ReadClassAvailabilityResponse, 0, len(scheduleCapacities))

	for _, cap := range scheduleCapacities {
		class, err := sc.ClassRepository.GetByID(tx, cap.ClassID)
		if err != nil || class == nil {
			return nil, fmt.Errorf("class not found for ID %d: %w", cap.ClassID, err)
		}

		occupied, err := sc.TicketRepository.CountByScheduleClassAndStatuses(tx, schedule.ID, cap.ClassID)
		if err != nil {
			return nil, fmt.Errorf("failed to count occupied tickets: %w", err)
		}

		manifest, err := sc.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, cap.ClassID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving manifest for ship %d, class %d: %w", schedule.ShipID, cap.ClassID, err)
		}
		if manifest == nil {
			log.Printf("skipping invalid/not found manifest for ship %d, class %d", schedule.ShipID, cap.ClassID)
			continue
		}

		fare, err := sc.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving fare for manifest %d, route %d: %w", manifest.ID, schedule.RouteID, err)
		}
		if fare == nil {
			log.Printf("skipping invalid/not found fare for manifest %d, route %d", manifest.ID, schedule.RouteID)
			continue
		}

		classAvailabilities = append(classAvailabilities, model.ReadClassAvailabilityResponse{
			ClassID:           cap.ClassID,
			ClassName:         class.ClassName,
			Type:              class.Type,
			TotalCapacity:     cap.Quota,
			AvailableCapacity: cap.Quota - int(occupied),
			Price:             fare.TicketPrice,
			Currency:          "IDR",
		})
	}

	return &model.ReadScheduleDetailsResponse{
		ScheduleID:          schedule.ID,
		RouteID:             schedule.RouteID,
		ClassesAvailability: classAvailabilities,
	}, nil
}

func (sc *ScheduleUsecase) CreateScheduleWithAllocation(ctx context.Context, request *model.WriteScheduleRequest) error {
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule := &domain.Schedule{
		RouteID:           request.RouteID,
		ShipID:            request.ShipID,
		DepartureDatetime: &request.DepartureDatetime,
		ArrivalDatetime:   &request.ArrivalDatetime,
		Status:            &request.Status,
	}

	if err := sc.ScheduleRepository.Create(tx, schedule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	ship, err := sc.ShipRepository.GetByID(tx, schedule.ShipID)
	if err != nil || ship == nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch ship %d: %w", schedule.ShipID, err)
	}

	manifests, err := sc.ManifestRepository.FindByShipID(tx, ship.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch manifests for ship %d: %w", ship.ID, err)
	}

	allocations := []*domain.Allocation{}
	var invalidManifests []string
	var missingFares []string

	for _, manifest := range manifests {
		if manifest.ClassID == 0 || manifest.Capacity <= 0 {
			invalidManifests = append(invalidManifests, fmt.Sprintf("ManifestID=%d, ShipID=%d", manifest.ID, manifest.ShipID))
			continue
		}

		fare, err := sc.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error retrieving fare for manifest %d: %w", manifest.ID, err)
		}
		if fare == nil {
			missingFares = append(missingFares, fmt.Sprintf("ManifestID=%d", manifest.ID))
			continue
		}

		allocation := &domain.Allocation{
			ScheduleID: schedule.ID,
			ClassID:    manifest.ClassID,
			Quota:      manifest.Capacity,
		}
		allocations = append(allocations, allocation)
	}

	if len(allocations) == 0 {
		tx.Rollback()
		var details []string
		if len(invalidManifests) > 0 {
			details = append(details, fmt.Sprintf("invalid manifests: [%s]", strings.Join(invalidManifests, ", ")))
		}
		if len(missingFares) > 0 {
			details = append(details, fmt.Sprintf("missing fares: [%s]", strings.Join(missingFares, ", ")))
		}
		return fmt.Errorf(
			"failed to create schedule for RouteID=%d, ShipID=%d. Details: %s",
			schedule.RouteID,
			schedule.ShipID,
			strings.Join(details, "; "),
		)
	}

	if err := sc.AllocationRepository.CreateBulk(tx, allocations); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create allocations: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
