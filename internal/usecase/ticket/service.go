package ticket

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type TicketUsecase struct {
	DB                 *gorm.DB
	TicketRepository   TicketRepository
	ScheduleRepository ScheduleRepository
	ManifestRepository ManifestRepository // Assuming ManifestRepository is defined elsewhere
	FareRepository     FareRepository     // Assuming FareRepository is defined elsewhere
}

func NewTicketUsecase(
	db *gorm.DB,
	ticket_repository TicketRepository,
	schedule_repository ScheduleRepository,
	manifest_repository ManifestRepository,
	fare_repository FareRepository,
) *TicketUsecase {
	return &TicketUsecase{
		DB:                 db,
		ScheduleRepository: schedule_repository,
		TicketRepository:   ticket_repository,
		ManifestRepository: nil, // Initialize if needed
		FareRepository:     nil, // Initialize if needed
	}
}

func (t *TicketUsecase) CreateTicket(ctx context.Context, request *model.WriteTicketRequest) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	schedule, err := t.ScheduleRepository.GetByID(tx, request.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule not found")
	}

	manifest, err := t.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, request.ClassID)
	if err != nil || manifest == nil {
		return fmt.Errorf("manifest missing for ship %d, class %d", schedule.ShipID, request.ClassID)
	}

	fare, err := t.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
	if err != nil || fare == nil {
		return fmt.Errorf("fare missing for manifest %d, route %d", manifest.ID, schedule.RouteID)
	}

	ticket := &entity.Ticket{
		ScheduleID:      request.ScheduleID,
		ClassID:         request.ClassID,
		BookingID:       request.BookingID,
		ClaimSessionID:  request.ClaimSessionID,
		Type:            request.Type,
		Price:           request.Price,
		PassengerName:   request.PassengerName,
		PassengerAge:    request.PassengerAge,
		PassengerGender: request.PassengerGender,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		SeatNumber:      request.SeatNumber,
		LicensePlate:    request.LicensePlate,
	}

	if err := t.TicketRepository.Create(tx, ticket); err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) GetAllTickets(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := t.TicketRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	tickets, err := t.TicketRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadTicketResponses(tickets), int(total), nil
}

func (t *TicketUsecase) GetAllTicketsByScheduleID(ctx context.Context, schedule_id, limit, offset int, sort, search string) ([]*model.ReadTicketResponse, int, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	var total int64
	countt := tx.Model(&entity.Ticket{}).Where("schedule_id = ?", schedule_id)
	if search != "" {
		search = "%" + search + "%"
		countt = countt.Where("passenger_name ILIKE ?", search)
	}
	if err := countt.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	tickets, err := t.TicketRepository.GetByScheduleID(tx, schedule_id, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadTicketResponses(tickets), int(total), nil
}

func (t *TicketUsecase) GetTicketByID(ctx context.Context, id uint) (*model.ReadTicketResponse, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	ticket, err := t.TicketRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}
	if ticket == nil {
		return nil, errors.New("ticket not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadTicketResponse(ticket), nil
}

func (t *TicketUsecase) UpdateTicket(ctx context.Context, request *model.UpdateTicketRequest) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	ticket, err := t.TicketRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}

	schedule, err := t.ScheduleRepository.GetByID(tx, request.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule not found")
	}

	manifest, err := t.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, request.ClassID)
	if err != nil || manifest == nil {
		return fmt.Errorf("manifest missing for ship %d, class %d", schedule.ShipID, request.ClassID)
	}

	fare, err := t.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
	if err != nil || fare == nil {
		return fmt.Errorf("fare missing for manifest %d, route %d", manifest.ID, schedule.RouteID)
	}

	ticket.ScheduleID = request.ScheduleID
	ticket.ClassID = request.ClassID
	ticket.BookingID = &request.BookingID
	ticket.ClaimSessionID = &request.ClaimSessionID
	ticket.Type = request.Type
	ticket.Price = fare.TicketPrice
	ticket.PassengerName = &request.PassengerName
	ticket.PassengerAge = &request.PassengerAge
	ticket.PassengerGender = &request.PassengerGender
	ticket.IDType = &request.IDType
	ticket.IDNumber = &request.IDNumber
	ticket.SeatNumber = request.SeatNumber
	ticket.LicensePlate = request.LicensePlate

	if err := t.TicketRepository.Update(tx, ticket); err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) DeleteTicket(ctx context.Context, id uint) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	ticket, err := t.TicketRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}

	if err := t.TicketRepository.Delete(tx, ticket); err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *TicketUsecase) CheckIn(ctx context.Context, id uint) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	if err := t.TicketRepository.CheckIn(tx, id); err != nil {
		return fmt.Errorf("failed to check-in ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
