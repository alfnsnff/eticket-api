package schedule

import (
	"context"
	"errors"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	Tx                   *tx.TxManager
	DB                   *gorm.DB // Assuming you have a DB field for the transaction manager
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
	tx *tx.TxManager,
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
		Tx:                   tx,
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
	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if request.DepartureDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if request.ArrivalDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if request.DepartureDatetime == request.ArrivalDatetime {
		return errors.New("schedule datetime cannot be same")
	}

	schedule := &entity.Schedule{
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

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ScheduleMapper.ToModels(schedules), int(total), nil
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

	return mapper.ScheduleMapper.ToModel(schedule), nil
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

	// Validate input
	if request.ID == 0 {
		return fmt.Errorf("allocation ID cannot be zero")
	}

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

	if schedule.ID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}

	if schedule.DepartureDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if schedule.ArrivalDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if schedule.DepartureDatetime == schedule.ArrivalDatetime {
		return errors.New("schedule datetime cannot be same")
	}

	return sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return sc.ScheduleRepository.Update(tx, schedule)
	})
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

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ScheduleMapper.ToModels(schedules), nil
}
func (sc *ScheduleUsecase) GetScheduleAvailability(ctx context.Context, scheduleID uint) (*model.ReadScheduleDetailsWithAvailabilityResponse, error) {
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

	classAvailabilities := make([]model.ScheduleClassAvailability, 0, len(scheduleCapacities))

	for _, cap := range scheduleCapacities {
		class, err := sc.ClassRepository.GetByID(tx, cap.ClassID)
		if err != nil || class == nil {
			return nil, fmt.Errorf("class not found for ID %d: %w", cap.ClassID, err)
		}

		occupied, err := sc.TicketRepository.CountByScheduleClassAndStatuses(tx, schedule.ID, cap.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"})
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

		classAvailabilities = append(classAvailabilities, model.ScheduleClassAvailability{
			ClassID:           cap.ClassID,
			ClassName:         class.ClassName,
			Type:              class.Type,
			TotalCapacity:     cap.Quota,
			AvailableCapacity: cap.Quota - int(occupied),
			Price:             fare.TicketPrice,
			Currency:          "IDR",
		})
	}

	return &model.ReadScheduleDetailsWithAvailabilityResponse{
		ScheduleID:          schedule.ID,
		RouteID:             schedule.RouteID,
		ClassesAvailability: classAvailabilities,
	}, nil
}
func (sc *ScheduleUsecase) CreateScheduleWithAllocation(ctx context.Context, request *model.WriteScheduleRequest) error {
	if request.DepartureDatetime.IsZero() {
		return errors.New("departure datetime cannot be empty")
	}
	if request.ArrivalDatetime.IsZero() {
		return errors.New("arrival datetime cannot be empty")
	}
	if request.ShipID == 0 {
		return errors.New("schedule ship ID cannot be zero")
	}

	schedule := mapper.ScheduleMapper.FromWrite(request)

	tx := sc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

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

	allocations := []*entity.Allocation{}
	for _, manifest := range manifests {
		if manifest.ClassID == 0 || manifest.Capacity <= 0 {
			log.Printf("Skipping invalid manifest %d for ship %d", manifest.ID, manifest.ShipID)
			continue
		}

		fare, err := sc.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error retrieving fare for manifest %d, route %d: %w", manifest.ID, schedule.RouteID, err)
		}
		if fare == nil {
			log.Printf("Skipping missing fare for manifest %d and route %d", manifest.ID, schedule.RouteID)
			continue
		}

		allocation := &entity.Allocation{
			ScheduleID: schedule.ID,
			ClassID:    manifest.ClassID,
			Quota:      manifest.Capacity,
		}
		allocations = append(allocations, allocation)
	}

	if len(allocations) == 0 {
		tx.Rollback()
		return fmt.Errorf("no valid manifest entries found for ship %d", ship.ID)
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
