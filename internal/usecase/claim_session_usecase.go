package usecase

import (
	"context"
	"errors"
	enum "eticket-api/internal/common/enums"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClaimSessionUsecase struct {
	DB                     *gorm.DB
	ClaimSessionRepository domain.ClaimSessionRepository
	ClaimItemRepository    domain.ClaimItemRepository // Assuming you have a ClaimItemRepository
	TicketRepository       domain.TicketRepository
	ScheduleRepository     domain.ScheduleRepository
	BookingRepository      domain.BookingRepository
	QuotaRepository        domain.QuotaRepository
}

func NewClaimSessionUsecase(
	db *gorm.DB,
	claim_session_repository domain.ClaimSessionRepository,
	claim_item_repository domain.ClaimItemRepository, // Assuming you have a ClaimItemRepository
	ticket_repository domain.TicketRepository,
	schedule_repository domain.ScheduleRepository,
	booking_repository domain.BookingRepository,
	quota_repository domain.QuotaRepository,
) *ClaimSessionUsecase {
	return &ClaimSessionUsecase{
		DB:                     db,
		ClaimSessionRepository: claim_session_repository,
		ClaimItemRepository:    claim_item_repository, // Assuming you have a ClaimItemRepository
		TicketRepository:       ticket_repository,
		ScheduleRepository:     schedule_repository,
		BookingRepository:      booking_repository,
		QuotaRepository:        quota_repository,
	}
}

func (cs *ClaimSessionUsecase) TESTCreateClaimSession(ctx context.Context, request *model.TESTWriteClaimSessionLockTicketsRequest) (*model.TESTReadClaimSessionLockTicketsResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Step 1: Validate schedule
	schedule, err := cs.ScheduleRepository.FindByID(tx, request.ScheduleID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		tx.Rollback()
		return nil, errs.ErrNotFound
	}

	// Step 2: Validate and lock quota for each class
	for _, item := range request.Items {
		// Get the quota
		quota, err := cs.QuotaRepository.FindByScheduleIDAndClassID(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to get quota for class %d: %w", item.ClassID, err)
		}
		if quota == nil {
			tx.Rollback()
			return nil, fmt.Errorf("quota not found for class %d", item.ClassID)
		}

		// Count current reserved
		used, err := cs.ClaimItemRepository.CountActiveReservedQuantity(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to count used quota: %w", err)
		}

		if used+int64(item.Quantity) > int64(quota.Quota) {
			tx.Rollback()
			return nil, fmt.Errorf("quota exceeded for class %d", item.ClassID)
		}
	}

	// Step 3: Create claim session
	claimSession := &domain.ClaimSession{
		SessionID:  uuid.NewString(),
		ScheduleID: request.ScheduleID,
		Status:     enum.ClaimSessionPendingData.String(),
		ExpiresAt:  time.Now().Add(13 * time.Minute),
	}
	if err := cs.ClaimSessionRepository.Insert(tx, claimSession); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create claim session: %w", err)
	}

	// Step 4: Create claim items
	for _, item := range request.Items {
		claimItem := &domain.ClaimItem{
			ClaimSessionID: claimSession.ID,
			ClassID:        item.ClassID,
			Quantity:       item.Quantity,
		}
		if err := cs.ClaimItemRepository.Insert(tx, claimItem); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create claim item: %w", err)
		}
	}

	// Step 5: Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Step 6: Return response
	return &model.TESTReadClaimSessionLockTicketsResponse{
		SessionID:  claimSession.SessionID,
		ExpiresAt:  claimSession.ExpiresAt,
		ClaimItems: request.Items, // You can map to another response struct if needed
	}, nil
}

func (cs *ClaimSessionUsecase) CreateClaimSession(ctx context.Context, request *model.WriteClaimSessionLockTicketsRequest) (*model.ReadClaimSessionLockTicketsResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	schedule, err := cs.ScheduleRepository.FindByID(tx, request.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		return nil, errs.ErrNotFound
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

	var ticketsToBuild []*domain.Ticket
	for _, item := range request.Items {
		quota, err := cs.QuotaRepository.FindByScheduleIDAndClassID(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			return nil, fmt.Errorf("failed to get quota capacity: %w", err)
		}
		occupied, err := cs.TicketRepository.CountByScheduleIDAndClassIDWithStatus(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			return nil, fmt.Errorf("failed to count tickets: %w", err)
		}

		available := int64(quota.Quota) - occupied
		if available < int64(item.Quantity) {
			return nil, fmt.Errorf("not enough slots for class %d (Available: %d, Requested: %d)", item.ClassID, available, item.Quantity)
		}

		for i := 0; i < int(item.Quantity); i++ {
			fmt.Println("Type:", quota.Class.Type, "Price:", quota.Price)
			ticketsToBuild = append(ticketsToBuild, &domain.Ticket{
				ScheduleID:     request.ScheduleID,
				ClassID:        item.ClassID,
				Price:          quota.Price,
				Type:           quota.Class.Type,
				ClaimSessionID: &claimSession.ID,
				IsCheckedIn:    false, // Default value
			})
		}
	}

	if err := cs.ClaimSessionRepository.Insert(tx, claimSession); err != nil {
		return nil, fmt.Errorf("failed to create claim session: %w", err)
	}

	if err := cs.TicketRepository.InsertBulk(tx, ticketsToBuild); err != nil {
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

func (cs *ClaimSessionUsecase) ListClaimSessions(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadClaimSessionResponse, int, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Count total claim sessions
	total, err := cs.ClaimSessionRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count claim sessions: %w", err)
	}

	// Retrieve claim sessions
	claimSessions, err := cs.ClaimSessionRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all claim sessions: %w", err)
	}

	responses := make([]*model.ReadClaimSessionResponse, len(claimSessions))
	for i, claimSession := range claimSessions {
		responses[i] = mapper.ClaimSessionToResponse(claimSession)
	}

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
		}
	}()

	session, err := cs.ClaimSessionRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ClaimSessionToResponse(session), nil
}

func (cs *ClaimSessionUsecase) GetBySessionID(ctx context.Context, sessionUUID string) (*model.ReadClaimSessionResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	session, err := cs.ClaimSessionRepository.FindBySessionID(tx, sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ClaimSessionToResponse(session), nil
}

func (cs *ClaimSessionUsecase) UpdateClaimSession(ctx context.Context, request *model.WriteClaimSessionDataEntryRequest, sessionID string) (*model.ReadClaimSessionDataEntryResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	session, err := cs.ClaimSessionRepository.FindBySessionID(tx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get claim session failed: %w", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("claim session expired")
	}
	if session == nil {
		return nil, errs.ErrNotFound
	}

	if session.Status != enum.ClaimSessionPendingData.String() {
		return nil, fmt.Errorf("claim session has invalid status: %s", session.Status)
	}

	// Build passenger map
	datas := make(map[uint]model.ClaimSessionTicketDataEntry)
	for _, data := range request.TicketData {
		datas[data.TicketID] = data
	}

	tickets, err := cs.TicketRepository.FindByClaimSessionID(tx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("retrieve tickets failed: %w", err)
	}

	schedule, err := cs.ScheduleRepository.FindByID(tx, session.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("retrieve schedule failed: %w", err)
	}

	orderID := utils.GenerateOrderID(
		schedule.DepartureHarbor.HarborAlias,
		schedule.ArrivalHarbor.HarborAlias,
		schedule.Ship.ShipAlias,
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

	if err := cs.BookingRepository.Insert(tx, booking); err != nil {
		return nil, fmt.Errorf("create booking failed: %w", err)
	}

	ticketsIds := make(map[uint]*domain.Ticket)
	for _, t := range tickets {
		if t.ID == 0 {
			continue
		}
		ticketsIds[t.ID] = t
	}

	var ticketToUpdate []*domain.Ticket
	var updatedTicketIDs []uint
	var total float64
	for id, data := range datas {
		ticket, ok := ticketsIds[id]
		if !ok {
			return nil, fmt.Errorf("ticket %d not found in session", id)
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

			count, err := cs.TicketRepository.CountByScheduleIDAndClassIDWithStatus(
				tx, ticket.ScheduleID, ticket.ClassID,
			)
			if err != nil {
				return nil, fmt.Errorf("seat number generation failed: %w", err)
			}
			seat := fmt.Sprintf("%s%d", ticket.Class.ClassAlias, count+1)
			ticket.SeatNumber = &seat

		case "vehicle":
			if data.LicensePlate == nil || *data.LicensePlate == "" {
				return nil, fmt.Errorf("missing license plate for vehicle ticket %d", id)
			}
			ticket.LicensePlate = data.LicensePlate
			ticket.SeatNumber = nil

		default:
			return nil, fmt.Errorf("unsupported ticket type for ticket %d - %s", id, ticket.Type)
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
		}
	}()

	session, err := cs.ClaimSessionRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return errs.ErrNotFound
	}

	if err := cs.ClaimSessionRepository.Delete(tx, session); err != nil {
		return fmt.Errorf("failed to delete fare: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
