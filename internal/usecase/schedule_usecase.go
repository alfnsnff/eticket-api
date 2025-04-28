package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ScheduleUsecase struct {
	DB                 *gorm.DB
	ScheduleRepository *repository.ScheduleRepository
	RouteRepository    *repository.RouteRepository
	PriceRepository    *repository.PriceRepository
	TicketRepository   *repository.TicketRepository
}

func NewScheduleUsecase(
	db *gorm.DB,
	scheduleRepository *repository.ScheduleRepository,
	routeRepository *repository.RouteRepository,
	priceRepository *repository.PriceRepository,
	ticketRepository *repository.TicketRepository,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		DB:                 db,
		ScheduleRepository: scheduleRepository,
		RouteRepository:    routeRepository,
		PriceRepository:    priceRepository,
		TicketRepository:   ticketRepository,
	}
}

// CreateSchedule validates and creates a new schedule
func (s *ScheduleUsecase) CreateSchedule(ctx context.Context, schedule *entities.Schedule) error {
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ScheduleRepository.Create(txDB, schedule)
	})
}

// GetPricesWithQuotaBySchedule retrieves ticket prices and remaining quota for a schedule
func (s *ScheduleUsecase) GetPricesWithQuotaBySchedule(ctx context.Context, scheduleID uint) ([]dto.ScheduleQuotaResponse, error) {
	var results []dto.ScheduleQuotaResponse

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(txDB, scheduleID)
		if err != nil {
			return err
		}
		if schedule == nil {
			return fmt.Errorf("schedule with ID %d not found", scheduleID)
		}

		prices, err := s.PriceRepository.GetByRouteID(txDB, schedule.RouteID)
		if err != nil {
			return err
		}

		for _, price := range prices {
			booked, err := s.TicketRepository.GetBookedCount(txDB, scheduleID, price.ID)
			if err != nil {
				return err
			}

			available := price.ShipClass.Capacity - booked

			results = append(results, dto.ScheduleQuotaResponse{
				PriceID:   price.ID,
				ClassName: price.ShipClass.Class.Name,
				Price:     price.Price,
				Capacity:  price.ShipClass.Capacity,
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
func (s *ScheduleUsecase) GetAllSchedules(ctx context.Context) ([]*entities.Schedule, error) {
	var schedules []*entities.Schedule

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		schedules, err = s.ScheduleRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	return schedules, nil
}

// GetScheduleByID retrieves a schedule by its ID
func (s *ScheduleUsecase) GetScheduleByID(ctx context.Context, id uint) (*entities.Schedule, error) {
	var schedule *entities.Schedule

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		schedule, err = s.ScheduleRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}

	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	return schedule, nil
}

// SearchSchedule searches a schedule by departure, arrival, and date
func (s *ScheduleUsecase) SearchSchedule(ctx context.Context, req dto.ScheduleSearchRequest) (*entities.Schedule, error) {
	var schedule *entities.Schedule

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		route, err := s.RouteRepository.Search(txDB, req.DepartureHarborID, req.ArrivalHarborID)
		if err != nil {
			return err
		}

		schedule, err = s.ScheduleRepository.Search(txDB, route.ID, req.Date, req.ShipID)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search schedule: %w", err)
	}

	return schedule, nil
}

// UpdateSchedule updates an existing schedule
func (s *ScheduleUsecase) UpdateSchedule(ctx context.Context, id uint, schedule *entities.Schedule) error {
	if id == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}

	schedule.ID = id

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ScheduleRepository.Update(txDB, schedule)
	})
}

// DeleteSchedule deletes a schedule by its ID
func (s *ScheduleUsecase) DeleteSchedule(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		schedule, err := s.ScheduleRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}
		return s.ScheduleRepository.Delete(txDB, id)
	})
}
