package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"eticket-api/pkg/utils/tx"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	Tx                   tx.TxManager
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
	tx tx.TxManager,
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

	if schedule.DepartureDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if schedule.ArrivalDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if schedule.DepartureDatetime == schedule.ArrivalDatetime {
		return errors.New("schedule datetime cannot be same")
	}

	// schedule.Status = "schedulled" // Set default status

	return sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return sc.ScheduleRepository.Create(tx, schedule)
	})
}

func (sc *ScheduleUsecase) GetAllSchedules(ctx context.Context, limit, offset int) ([]*model.ReadScheduleResponse, int, error) {
	schedules := []*entity.Schedule{}
	var total int64
	err := sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = sc.ScheduleRepository.Count(tx)
		if err != nil {
			return err
		}
		schedules, err = sc.ScheduleRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all schedules: %w", err)
	}

	return mapper.ScheduleMapper.ToModels(schedules), int(total), nil
}

func (sc *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	schedule := new(entity.Schedule)

	err := sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

	return sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

func (sc *ScheduleUsecase) GetAllScheduled(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	schedules := []*entity.Schedule{}

	err := sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		schedules, err = sc.ScheduleRepository.GetActiveSchedule(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	return mapper.ScheduleMapper.ToModels(schedules), nil
}

func (sc *ScheduleUsecase) GetScheduleAvailability(ctx context.Context, scheduleID uint) (*model.ReadScheduleDetailsWithAvailabilityResponse, error) {
	var schedule *entity.Schedule
	var classAvailabilities []model.ScheduleClassAvailability

	err := sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error

		schedule, err = sc.HelperGetSchedule(tx, scheduleID)
		if err != nil {
			return err
		}

		classAvailabilities, err = sc.HelperGetAvailability(tx, schedule)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get schedule details with availability: %w", err)
	}

	return &model.ReadScheduleDetailsWithAvailabilityResponse{
		ScheduleID:          schedule.ID,
		RouteID:             schedule.RouteID,
		ClassesAvailability: classAvailabilities,
	}, nil
}

func (sc *ScheduleUsecase) HelperGetSchedule(tx *gorm.DB, scheduleID uint) (*entity.Schedule, error) {
	schedule, err := sc.ScheduleRepository.GetByID(tx, scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	return schedule, nil
}

func (sc *ScheduleUsecase) HelperGetAvailability(tx *gorm.DB, schedule *entity.Schedule) ([]model.ScheduleClassAvailability, error) {
	scheduleCapacities, err := sc.AllocationRepository.FindByScheduleID(tx, schedule.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule capacities: %w", err)
	}

	result := make([]model.ScheduleClassAvailability, 0, len(scheduleCapacities))

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
		if err != nil || manifest == nil {
			return nil, fmt.Errorf("manifest not found for ship %d, class %d: %w", schedule.ShipID, cap.ClassID, err)
		}

		fare, err := sc.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
		if err != nil || fare == nil {
			return nil, fmt.Errorf("fare not found for manifest %d, route %d: %w", manifest.ID, schedule.RouteID, err)
		}

		result = append(result, model.ScheduleClassAvailability{
			ClassID:           cap.ClassID,
			ClassName:         class.ClassName,
			TotalCapacity:     cap.Quota,
			AvailableCapacity: cap.Quota - int(occupied),
			Price:             fare.TicketPrice,
			Currency:          "IDR",
		})
	}

	return result, nil
}

func (sc *ScheduleUsecase) CreateScheduleWithAllocation(ctx context.Context, request *model.WriteScheduleRequest) error {
	schedule := mapper.ScheduleMapper.FromWrite(request)

	if err := HelperValidateScheduleInput(schedule); err != nil {
		return err
	}

	err := sc.Tx.Execute(ctx, func(tx *gorm.DB) error {
		if err := sc.ScheduleRepository.Create(tx, schedule); err != nil {
			return fmt.Errorf("failed to create schedule: %w", err)
		}

		ship, err := sc.ShipRepository.GetByID(tx, schedule.ShipID)
		if err != nil || ship == nil {
			return fmt.Errorf("failed to fetch ship %d: %w", schedule.ShipID, err)
		}

		manifests, err := sc.ManifestRepository.FindByShipID(tx, ship.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch manifests for ship %d: %w", ship.ID, err)
		}

		allocations := HelperBuildAllocationsFromManifests(schedule.ID, manifests)
		if len(allocations) == 0 {
			return fmt.Errorf("no valid manifest entries found for ship %d", ship.ID)
		}

		if err := sc.AllocationRepository.CreateBulk(tx, allocations); err != nil {
			return fmt.Errorf("failed to create allocations: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("schedule creation transaction failed: %w", err)
	}

	return nil
}

func HelperValidateScheduleInput(schedule *entity.Schedule) error {
	if schedule.DepartureDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if schedule.ArrivalDatetime.IsZero() {
		return errors.New("schedule datetime cannot be empty")
	}
	if schedule.DepartureDatetime == schedule.ArrivalDatetime {
		return errors.New("schedule datetime cannot be same")
	}
	if schedule.ShipID == 0 {
		return errors.New("schedule ship ID cannot be zero")
	}
	// Add RouteID or additional checks here
	return nil
}

func HelperBuildAllocationsFromManifests(scheduleID uint, manifests []*entity.Manifest) []*entity.Allocation {
	allocations := []*entity.Allocation{}
	for _, manifest := range manifests {
		if manifest.ClassID == 0 || manifest.Capacity <= 0 {
			log.Printf("Skipping invalid manifest %d for ship %d", manifest.ID, manifest.ShipID)
			continue
		}
		allocation := &entity.Allocation{
			ScheduleID: scheduleID,
			ClassID:    manifest.ClassID,
			Quota:      manifest.Capacity,
		}
		allocations = append(allocations, allocation)
	}
	return allocations
}
