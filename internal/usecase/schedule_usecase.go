package usecase

import (
	"errors"
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	"fmt"
)

type ScheduleUsecase struct {
	ScheduleRepository *repository.ScheduleRepository
	RouteRepository    *repository.RouteRepository
	PriceRepository    *repository.PriceRepository
	TicketRepository   *repository.TicketRepository
}

func NewScheduleUsecase(
	scheduleRepository *repository.ScheduleRepository,
	routeRepository *repository.RouteRepository,
	priceRepository *repository.PriceRepository,
	ticketRepository *repository.TicketRepository) ScheduleUsecase {
	return ScheduleUsecase{
		ScheduleRepository: scheduleRepository,
		RouteRepository:    routeRepository,
		PriceRepository:    priceRepository,
		TicketRepository:   ticketRepository}
}

// Createschedule validates and creates a new schedule
func (s *ScheduleUsecase) CreateSchedule(schedule *entities.Schedule) error {
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}
	return s.ScheduleRepository.Create(schedule)
}

// TicketUsecase.go
func (s *ScheduleUsecase) GetPricesWithQuotaBySchedule(scheduleID uint) ([]dto.ScheduleQuotaResponse, error) {
	schedule, err := s.ScheduleRepository.GetByID(scheduleID)
	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, fmt.Errorf("schedule with ID %d not found", scheduleID)
	}

	if s.PriceRepository == nil {
		return nil, fmt.Errorf("PriceRepository is nil")
	}

	prices, err := s.PriceRepository.GetByRouteID(1)

	if err != nil {
		return nil, err
	}

	var results []dto.ScheduleQuotaResponse
	for _, price := range prices {
		booked, err := s.TicketRepository.GetBookedCount(scheduleID, price.ID)
		if err != nil {
			return nil, err
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

	return results, nil
}

// GetAllschedulees retrieves all schedules
func (s *ScheduleUsecase) GetAllSchedules() ([]*entities.Schedule, error) {
	return s.ScheduleRepository.GetAll()
}

// GetscheduleByID retrieves a schedule by its ID
func (s *ScheduleUsecase) GetScheduleByID(id uint) (*entities.Schedule, error) {
	schedule, err := s.ScheduleRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("schedule not found")
	}
	return schedule, nil
}

func (s *ScheduleUsecase) SearchSchedule(req dto.ScheduleSearchRequest) (*entities.Schedule, error) {
	route, err := s.RouteRepository.Search(req.DepartureHarborID, req.ArrivalHarborID)

	if err != nil {
		return nil, err
	}

	return s.ScheduleRepository.Search(route.ID, req.Date, req.ShipID)
}

// Updateschedule updates an existing schedule
func (s *ScheduleUsecase) UpdateSchedule(id uint, schedule *entities.Schedule) error {
	schedule.ID = id

	if schedule.ID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}
	if schedule.Datetime.IsZero() {
		return fmt.Errorf("schedule datetime cannot be empty")
	}
	return s.ScheduleRepository.Update(schedule)
}

// Deleteschedule deletes a schedule by its ID
func (s *ScheduleUsecase) DeleteSchedule(id uint) error {
	schedule, err := s.ScheduleRepository.GetByID(id)
	if err != nil {
		return err
	}
	if schedule == nil {
		return errors.New("schedule not found")
	}
	return s.ScheduleRepository.Delete(id)
}
