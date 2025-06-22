package claim_session

import (
	"context"
	"errors"
	enum "eticket-api/internal/common/enums"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClaimSessionUsecase struct {
	DB                     *gorm.DB
	ClaimSessionRepository ClaimSessionRepository
	TicketRepository       TicketRepository
	ScheduleRepository     ScheduleRepository
	AllocationRepository   AllocationRepository // Your AllocationRepository implements this
	ManifestRepository     ManifestRepository
	FareRepository         FareRepository
	BookingRepository      BookingRepository
}

func NewClaimSessionUsecase(
	db *gorm.DB,
	claim_session_repository ClaimSessionRepository,
	ticket_repository TicketRepository,
	schedule_repository ScheduleRepository,
	allocation_repository AllocationRepository,
	manifest_repository ManifestRepository,
	fare_repository FareRepository,
	booking_repository BookingRepository,
) *ClaimSessionUsecase {
	return &ClaimSessionUsecase{
		DB:                     db,
		ClaimSessionRepository: claim_session_repository,
		TicketRepository:       ticket_repository,
		ScheduleRepository:     schedule_repository,
		AllocationRepository:   allocation_repository,
		ManifestRepository:     manifest_repository,
		FareRepository:         fare_repository,
		BookingRepository:      booking_repository,
	}
}

func (cs *ClaimSessionUsecase) CreateClaimSession(ctx context.Context, request *model.WriteClaimSessionLockTicketsRequest) (*model.ReadClaimSessionLockTicketsResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	for id, item := range request.Items {
		if item.Quantity == 0 {
			return nil, fmt.Errorf("missing quantity field for item %d", id)
		}
		if item.ClassID == 0 {
			return nil, fmt.Errorf("missing class field for item %d", id)
		}
	}

	schedule, err := cs.ScheduleRepository.GetByID(tx, request.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("schedule not found")
	}

	for _, item := range request.Items {
		cap, err := cs.AllocationRepository.LockByScheduleAndClass(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			return nil, fmt.Errorf("failed to lock capacity: %w", err)
		}
		if cap == nil {
			return nil, fmt.Errorf("allocation not found for class %d schedule %d", item.ClassID, request.ScheduleID)
		}

		count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			return nil, fmt.Errorf("failed to count tickets: %w", err)
		}

		available := int64(cap.Quota) - count
		if available < int64(item.Quantity) {
			return nil, fmt.Errorf("not enough slots for class %d (Available: %d, Requested: %d)", item.ClassID, available, item.Quantity)
		}
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session UUID: %w", err)
	}

	claimSession := &domain.ClaimSession{
		SessionID:  uuid.String(),
		ScheduleID: request.ScheduleID,
		Status:     enum.ClaimSessionPendingData.String(),
		ExpiresAt:  time.Now().Add(13 * time.Minute),
	}

	if err := cs.ClaimSessionRepository.Create(tx, claimSession); err != nil {
		return nil, fmt.Errorf("failed to create claim session: %w", err)
	}

	var ticketsToBuild []*domain.Ticket

	for _, item := range request.Items {
		manifest, err := cs.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, item.ClassID)
		if err != nil || manifest == nil {
			return nil, fmt.Errorf("manifest missing for ship %d, class %d", schedule.ShipID, item.ClassID)
		}

		fare, err := cs.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
		if err != nil || fare == nil {
			return nil, fmt.Errorf("fare missing for manifest %d, route %d", manifest.ID, schedule.RouteID)
		}

		for i := 0; i < int(item.Quantity); i++ {

			ticketsToBuild = append(ticketsToBuild, &domain.Ticket{
				ScheduleID:     request.ScheduleID,
				ClassID:        item.ClassID,
				Price:          fare.TicketPrice,
				Type:           manifest.Class.Type,
				ClaimSessionID: &claimSession.ID,
			})
		}
	}

	if err := cs.TicketRepository.CreateBulk(tx, ticketsToBuild); err != nil {
		return nil, fmt.Errorf("failed to create tickets: %w", err)
	}

	var claimedTicketIDs []uint
	for _, ticket := range ticketsToBuild {
		claimedTicketIDs = append(claimedTicketIDs, ticket.ID)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.ReadClaimSessionLockTicketsResponse{
		ClaimedTicketIDs: claimedTicketIDs,
		ExpiresAt:        claimSession.ExpiresAt,
		SessionID:        claimSession.SessionID,
	}, nil
}

func (cs *ClaimSessionUsecase) GetAllClaimSessions(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadClaimSessionResponse, int, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Count total claim sessions
	total, err := cs.ClaimSessionRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count claim sessions: %w", err)
	}

	// Retrieve claim sessions
	claimSessions, err := cs.ClaimSessionRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all claim sessions: %w", err)
	}

	responses := make([]*model.ReadClaimSessionResponse, len(claimSessions))
	for i, claimSession := range claimSessions {
		responses[i] = ClaimSessionToResponse(claimSession)
	}

	// // Map claim sessions to response models
	// var responses []*model.ReadClaimSessionResponse
	// for _, session := range claimSessions {
	// 	tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
	// 	if err != nil {
	// 		return nil, 0, fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
	// 	}

	// 	responses = append(responses, ToReadClaimSessionResponse(session, tickets))
	// }

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (cs *ClaimSessionUsecase) GetClaimSessionByID(ctx context.Context, id uint) (*model.ReadClaimSessionResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	session, err := cs.ClaimSessionRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ClaimSessionToResponse(session), nil
}

func (cs *ClaimSessionUsecase) GetBySessionID(ctx context.Context, sessionUUID string) (*model.ReadClaimSessionResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if sessionUUID == "" {
		return nil, errors.New("invalid request: SessionID is required")
	}

	session, err := cs.ClaimSessionRepository.GetByUUID(tx, sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ClaimSessionToResponse(session), nil
}

func (cs *ClaimSessionUsecase) UpdateClaimSession(ctx context.Context, request *model.WriteClaimSessionDataEntryRequest, sessionID string) (*model.ReadClaimSessionDataEntryResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	session, err := cs.ClaimSessionRepository.GetByUUIDWithLock(tx, sessionID, true)
	if err != nil {
		return nil, fmt.Errorf("get claim session failed: %w", err)
	}

	if session == nil || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("claim session not found or expired")
	}

	if session.Status != enum.ClaimSessionPendingData.String() {
		return nil, fmt.Errorf("claim session has invalid status: %s", session.Status)
	}

	// Build passenger map
	datas := make(map[uint]model.ClaimSessionTicketDataEntry)
	for _, data := range request.TicketData {
		datas[data.TicketID] = data
	}

	tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("retrieve tickets failed: %w", err)
	}

	schedule, err := cs.ScheduleRepository.GetByID(tx, session.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("retrieve schedule failed: %w", err)
	}

	orderID := utils.GenerateOrderID(
		*schedule.Route.DepartureHarbor.HarborAlias,
		*schedule.Route.ArrivalHarbor.HarborAlias,
		*schedule.Ship.ShipAlias,
		time.Now(),
	)

	booking := &domain.Booking{
		OrderID:      &orderID,
		ScheduleID:   session.ScheduleID,
		IDType:       request.IDType,
		IDNumber:     request.IDNumber,
		PhoneNumber:  request.PhoneNumber,
		CustomerName: request.CustomerName,
		Email:        request.Email,
	}

	if err := cs.BookingRepository.Create(tx, booking); err != nil {
		return nil, fmt.Errorf("create booking failed: %w", err)
	}

	ticketsIds := make(map[uint]*domain.Ticket)
	for _, t := range tickets {
		if t.ID == 0 {
			continue
		}
		ticketsIds[t.ID] = t // ✅ correct: assigning to the actual map
	}

	var ticketToUpdate []*domain.Ticket
	var updatedTicketIDs []uint
	var total float32
	for id, data := range datas {
		ticket, ok := ticketsIds[id]
		if !ok {
			return nil, fmt.Errorf("ticket %d not found in session", id)
		}

		if data.PassengerName == "" || data.PassengerAge == 0 || data.Address == "" {
			return nil, fmt.Errorf("missing passenger data for ticket %d", id)
		}

		ticket.PassengerName = &data.PassengerName
		ticket.PassengerAge = &data.PassengerAge
		ticket.PassengerGender = &data.PassengerGender
		ticket.Address = &data.Address
		ticket.BookingID = &booking.ID

		switch ticket.Type {
		case "passenger":
			if data.IDType == "" || data.IDNumber == "" {
				return nil, fmt.Errorf("missing ID info for passenger ticket %d", id)
			}
			ticket.IDType = &data.IDType
			ticket.IDNumber = &data.IDNumber

			count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(
				tx, ticket.ScheduleID, ticket.ClassID,
			)
			if err != nil {
				return nil, fmt.Errorf("seat number generation failed: %w", err)
			}
			seat := fmt.Sprintf("%s%d", *ticket.Class.ClassAlias, count+1)
			ticket.SeatNumber = &seat

		case "vehicle":
			if data.LicensePlate == nil || *data.LicensePlate == "" {
				return nil, fmt.Errorf("missing license plate for vehicle ticket %d", id)
			}
			ticket.LicensePlate = data.LicensePlate
			ticket.SeatNumber = nil

		default:
			return nil, fmt.Errorf("unsupported ticket type for ticket %d", id)
		}

		total += ticket.Price
		ticketToUpdate = append(ticketToUpdate, ticket)
		updatedTicketIDs = append(updatedTicketIDs, ticket.ID)
	}

	if len(ticketToUpdate) == len(tickets) {
		if err := cs.TicketRepository.UpdateBulk(tx, ticketToUpdate); err != nil {
			return nil, fmt.Errorf("update tickets failed: %w", err)
		}
	}

	session = &domain.ClaimSession{
		ID:         session.ID,
		SessionID:  session.SessionID,
		ScheduleID: session.ScheduleID,
		Status:     enum.ClaimSessionPendingPayment.String(),
		ExpiresAt:  session.ExpiresAt.Add(8 * time.Minute), // Extend expiration for payment
	}

	if err := cs.ClaimSessionRepository.Update(tx, session); err != nil {
		return nil, fmt.Errorf("update claim session failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.ReadClaimSessionDataEntryResponse{
		BookingID:        booking.ID,
		OrderID:          orderID,
		UpdatedTicketIDs: updatedTicketIDs,
	}, nil
}

func (cs *ClaimSessionUsecase) DeleteClaimSession(ctx context.Context, id uint) error {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	session, err := cs.ClaimSessionRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return errors.New("session not found")
	}

	if err := cs.ClaimSessionRepository.Delete(tx, session); err != nil {
		return fmt.Errorf("failed to delete fare: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
