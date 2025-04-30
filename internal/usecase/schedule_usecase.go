package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	DB                 *gorm.DB
	ScheduleRepository *repository.ScheduleRepository
	RouteRepository    *repository.RouteRepository
	PriceRepository    *repository.FareRepository
	TicketRepository   *repository.TicketRepository
}

func NewScheduleUsecase(
	db *gorm.DB,
	schedule_repository *repository.ScheduleRepository,
	route_repository *repository.RouteRepository,
	fare_repository *repository.FareRepository,
	ticket_repository *repository.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		DB:                 db,
		ScheduleRepository: schedule_repository,
		RouteRepository:    route_repository,
		PriceRepository:    fare_repository,
		TicketRepository:   ticket_repository,
	}
}

// CreateSchedule validates and creates a new schedule
func (s *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.WriteScheduleRequest) error {
	schedule := mapper.ToScheduleEntity(request)
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ScheduleRepository.Create(tx, schedule)
	})
}

// GetPricesWithQuotaBySchedule retrieves ticket prices and remaining quota for a schedule
func (s *ScheduleUsecase) GetPricesWithQuotaBySchedule(ctx context.Context, scheduleID uint) ([]model.ScheduleQuotaResponse, error) {
	var results []model.ScheduleQuotaResponse

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(tx, scheduleID)
		if err != nil {
			return err
		}
		if schedule == nil {
			return fmt.Errorf("schedule with ID %d not found", scheduleID)
		}

		prices, err := s.PriceRepository.GetByRouteID(tx, schedule.RouteID)
		if err != nil {
			return err
		}

		for _, price := range prices {
			booked, err := s.TicketRepository.GetBookedCount(tx, scheduleID, price.ID)
			if err != nil {
				return err
			}

			available := price.Manifest.Capacity - booked

			results = append(results, model.ScheduleQuotaResponse{
				PriceID:   price.ID,
				ClassName: price.Manifest.Class.Name,
				Price:     price.Price,
				Capacity:  price.Manifest.Capacity,
				Booked:    booked,
				Available: available,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get prices and quotas: %w", err)
	}

	return results, nil
}

// GetAllSchedules retrieves all schedules
func (s *ScheduleUsecase) GetAllSchedules(ctx context.Context) ([]*model.ReadScheduleResponse, error) {
	schedules := []*entities.Schedule{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		schedules, err = s.ScheduleRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	return mapper.ToSchedulesModel(schedules), nil
}

// GetScheduleByID retrieves a schedule by its ID
func (s *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*model.ReadScheduleResponse, error) {
	schedule := new(entities.Schedule)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		schedule, err = s.ScheduleRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}

	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	return mapper.ToScheduleModel(schedule), nil
}

// SearchSchedule searches a schedule by departure, arrival, and date
func (s *ScheduleUsecase) SearchSchedule(ctx context.Context, request *model.SearchScheduleRequest) (*model.ReadScheduleResponse, error) {
	schedule := new(entities.Schedule)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		route, err := s.RouteRepository.Search(tx, request.DepartureHarborID, request.ArrivalHarborID)
		if err != nil {
			return err
		}

		schedule, err = s.ScheduleRepository.Search(tx, route.ID, request.Date, request.ShipID)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search schedule: %w", err)
	}

	return mapper.ToScheduleModel(schedule), nil
}

// UpdateSchedule updates an existing schedule
func (s *ScheduleUsecase) UpdateSchedule(ctx context.Context, id uint, request *model.WriteScheduleRequest) error {
	schedule := mapper.ToScheduleEntity(request)
	schedule.ID = id

	if schedule.ID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}

	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ScheduleRepository.Update(tx, schedule)
	})
}

// DeleteSchedule deletes a schedule by its ID
func (s *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}
		return s.ScheduleRepository.Delete(tx, schedule)
	})
}
