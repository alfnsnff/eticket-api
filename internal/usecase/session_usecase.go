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
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionUsecase struct {
	Tx                   *tx.TxManager
	SessionRepository    *repository.SessionRepository
	TicketRepository     *repository.TicketRepository
	ScheduleRepository   *repository.ScheduleRepository
	AllocationRepository *repository.AllocationRepository // Your AllocationRepository implements this
	ManifestRepository   *repository.ManifestRepository
	FareRepository       *repository.FareRepository
}

func NewSessionUsecase(
	tx *tx.TxManager,
	sessionRepo *repository.SessionRepository,
	ticketRepo *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	allocation_repository *repository.AllocationRepository,
	manifest_repository *repository.ManifestRepository,
	fare_repository *repository.FareRepository,
) *SessionUsecase {
	return &SessionUsecase{
		Tx:                   tx,
		SessionRepository:    sessionRepo,
		TicketRepository:     ticketRepo,
		ScheduleRepository:   schedule_repository,
		AllocationRepository: allocation_repository,
		ManifestRepository:   manifest_repository,
		FareRepository:       fare_repository,
	}
}

func (cs *SessionUsecase) CreateSession(ctx context.Context, request *model.WriteClaimSessionRequest) error {
	session := mapper.SessionMapper.FromWrite(request)

	return cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return cs.SessionRepository.Create(tx, session)
	})
}

func (cs *SessionUsecase) GetAllSessions(ctx context.Context, limit, offset int) ([]*model.ReadClaimSessionResponse, int, error) {
	sessions := []*entity.ClaimSession{}
	var total int64
	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = cs.SessionRepository.Count(tx)
		if err != nil {
			return err
		}
		sessions, err = cs.SessionRepository.GetAll(tx, limit, offset)
		return err
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all sessions: %w", err)
	}

	return cs.HelperBuildSessionListResponse(ctx, sessions), int(total), nil
}

func (cs *SessionUsecase) GetSessionByID(ctx context.Context, id uint) (*model.ReadClaimSessionResponse, error) {
	session := new(entity.ClaimSession)
	tickets := []*entity.Ticket{}

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		session, err = cs.SessionRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if session == nil {
			return errors.New("session not found")
		}

		tickets, err = cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cs.HelperBuildResponse(session, tickets), nil
}

func (cs *SessionUsecase) UpdateSession(ctx context.Context, id uint, request *model.UpdateClaimSessionRequest) error {
	session := mapper.SessionMapper.FromUpdate(request)
	session.ID = id

	if session.ID == 0 {
		return fmt.Errorf("session ID cannot be zero")
	}

	return cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return cs.SessionRepository.Update(tx, session)
	})
}

func (cs *SessionUsecase) DeleteSession(ctx context.Context, id uint) error {
	return cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		session, err := cs.SessionRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if session == nil {
			return errors.New("session not found")
		}
		return cs.SessionRepository.Delete(tx, session)
	})
}

func (cs *SessionUsecase) GetBySessionID(ctx context.Context, sessionUUID string) (*model.ReadClaimSessionResponse, error) {
	if sessionUUID == "" {
		return nil, errors.New("invalid request: SessionID is required")
	}

	session := new(entity.ClaimSession)
	tickets := []*entity.Ticket{}

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		session, err = cs.SessionRepository.GetByUUID(tx, sessionUUID)
		if err != nil {
			return err
		}
		if session == nil {
			return errors.New("session not found")
		}

		tickets, err = cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cs.HelperBuildResponse(session, tickets), nil
}

func (cs *SessionUsecase) SessionLockTickets(ctx context.Context, request *model.ClaimedSessionLockTicketsRequest) (*model.ClaimedSessionLockTicketsResponse, error) {
	if err := cs.HelperValidateLockRequest(request); err != nil {
		return nil, err
	}

	var claimedTicketIDs []uint
	var expiryTime time.Time
	var createdSessionUUID string // To hold the generated UUID

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		_, err := cs.HelperLockAndCheckAvailability(tx, request)
		if err != nil {
			return err
		}

		now := time.Now()
		expiryTime = time.Now().Add(15 * time.Minute)

		sessionUUID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("failed to generate session UUID: %w", err) // This will trigger rollback
		}
		createdSessionUUID = sessionUUID.String() // Store for response

		tickets, err := cs.HelperBuildTickets(tx, request)
		if err != nil {
			return err
		}

		newClaimSession := &entity.ClaimSession{
			SessionID:  createdSessionUUID,
			ScheduleID: request.ScheduleID,
			ClaimedAt:  now,
			ExpiresAt:  expiryTime,
			// Other fields like CreatedAt/UpdatedAt handled by GORM
		}

		err = cs.SessionRepository.Create(tx, newClaimSession) // Use txDB
		if err != nil {
			return fmt.Errorf("failed to create claim session: %w", err) // This will trigger rollback
		}

		// Link the newly created tickets to the ClaimSession
		for _, ticket := range tickets {
			ticket.ClaimSessionID = &newClaimSession.ID // Set the FK to the new ClaimSession ID
		}

		if err := cs.TicketRepository.CreateBulk(tx, tickets); err != nil {
			return fmt.Errorf("failed to create tickets: %w", err)
		}

		for _, t := range tickets {
			claimedTicketIDs = append(claimedTicketIDs, t.ID)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to claim tickets: %w", err)
	}

	return &model.ClaimedSessionLockTicketsResponse{
		ClaimedTicketIDs: claimedTicketIDs,
		ExpiresAt:        expiryTime,
		SessionID:        createdSessionUUID,
	}, nil
}

func (cs *SessionUsecase) SessionDataEntry(ctx context.Context, request *model.ClaimedSessionFillPassengerDataRequest, sessionID string) (*model.ClaimedSessionFillPassengerDataResponse, error) {
	if len(request.PassengerData) == 0 {
		return nil, errors.New("invalid request: UserID and passenger data are required")
	}

	_, passengerMap := HelperExtractPassengerData(request)

	var updatedIDs []uint
	var failed []model.ClaimedSessionTicketUpdateFailure

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		session, err := cs.SessionRepository.GetByUUIDWithLock(tx, sessionID, true)
		if err != nil {
			return fmt.Errorf("failed to retrieve claim session %s within transaction: %w", sessionID, err)
		}
		if session == nil {
			return errors.New("claim session not found")
		}

		now := time.Now()
		if session.ExpiresAt.Before(now) {
			return errors.New("claim session has expired")
		}

		tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets: %w", err)
		}

		updatedIDs, failed, tickets = cs.HelperValidateAndUpdateTickets(tickets, passengerMap, now)

		if len(tickets) > 0 {
			err = cs.TicketRepository.UpdateBulk(tx, tickets)
			if err != nil {
				return fmt.Errorf("failed to save tickets: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fill passenger data failed: %w", err)
	}

	return &model.ClaimedSessionFillPassengerDataResponse{
		UpdatedTicketIDs: updatedIDs,
		FailedTickets:    failed,
	}, nil
}

func (cs *SessionUsecase) HelperValidateLockRequest(request *model.ClaimedSessionLockTicketsRequest) error {
	if request.ScheduleID == 0 || len(request.Items) == 0 {
		return errors.New("invalid claim request")
	}
	return nil
}

func (cs *SessionUsecase) HelperLockAndCheckAvailability(tx *gorm.DB, request *model.ClaimedSessionLockTicketsRequest) (map[uint]int64, error) {
	checks := make(map[uint]int64)
	for _, item := range request.Items {
		if item.Quantity == 0 {
			continue
		}

		cap, err := cs.AllocationRepository.LockByScheduleAndClass(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			return nil, fmt.Errorf("failed to lock capacity: %w", err)
		}
		if cap == nil {
			return nil, fmt.Errorf("allocation not found for class %d", item.ClassID)
		}

		count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(tx, request.ScheduleID, item.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"})
		if err != nil {
			return nil, fmt.Errorf("failed to count tickets: %w", err)
		}

		available := int64(cap.Quota) - count
		if available < int64(item.Quantity) {
			return nil, fmt.Errorf("not enough slots for class %d (Available: %d, Requested: %d)", item.ClassID, available, item.Quantity)
		}

		checks[item.ClassID] = available
	}
	return checks, nil
}

func (cs *SessionUsecase) HelperBuildTickets(tx *gorm.DB, request *model.ClaimedSessionLockTicketsRequest) ([]*entity.Ticket, error) {
	schedule, err := cs.ScheduleRepository.GetByID(tx, request.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	var tickets []*entity.Ticket
	// now := time.Now()

	for _, item := range request.Items {
		if item.Quantity == 0 {
			continue
		}

		manifest, err := cs.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, item.ClassID)
		if err != nil || manifest == nil {
			return nil, fmt.Errorf("manifest missing for ship %d, class %d", schedule.ShipID, item.ClassID)
		}

		fare, err := cs.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
		if err != nil || fare == nil {
			return nil, fmt.Errorf("fare missing for manifest %d, route %d", manifest.ID, schedule.RouteID)
		}

		for i := 0; i < int(item.Quantity); i++ {
			tickets = append(tickets, &entity.Ticket{
				ScheduleID: request.ScheduleID,
				ClassID:    item.ClassID,
				Status:     "pending_data_entry",
				Price:      fare.TicketPrice,
				// ClaimedAt:  now,
				// ExpiresAt:  expiry,
				ClaimSessionID: nil, // Will be set after ClaimSession is created
			})
		}
	}
	return tickets, nil
}

// buildSessionResponse generates a consistent ReadClaimSessionResponse, optionally using ticket info.
func (cs *SessionUsecase) HelperBuildResponse(session *entity.ClaimSession, tickets []*entity.Ticket) *model.ReadClaimSessionResponse {
	var ticketPrices []model.ClaimSessionTicketPricesResponse
	var ticketDetails []model.ClaimedSessionTicketDetailResponse
	var total float32

	if len(tickets) > 0 {
		ticketPrices, total = cs.HelperBuildPriceBreakdown(tickets)
		ticketDetails = cs.HelperBuildTicketBreakdown(tickets)
	} else {
		ticketPrices = []model.ClaimSessionTicketPricesResponse{}
		total = 0
	}

	return &model.ReadClaimSessionResponse{
		ID:          session.ID,
		SessionID:   session.SessionID,
		ScheduleID:  session.ScheduleID,
		Schedule:    *mapper.ScheduleSessionMapper.ToModel(&session.Schedule),
		ClaimedAt:   session.ClaimedAt,
		ExpiresAt:   session.ExpiresAt,
		Prices:      ticketPrices,
		Tickets:     ticketDetails,
		TotalAmount: total,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}
}

// buildTicketPriceBreakdown groups tickets by class and calculates subtotals and total.
func (cs *SessionUsecase) HelperBuildTicketBreakdown(tickets []*entity.Ticket) []model.ClaimedSessionTicketDetailResponse {
	result := make([]model.ClaimedSessionTicketDetailResponse, len(tickets))
	for i, v := range tickets {
		result[i] = model.ClaimedSessionTicketDetailResponse{
			TicketID: v.ID,
			Class:    *mapper.TicketClassToSessionClassMapper.ToModel(&v.Class),
			Price:    v.Price,
		}
	}
	return result
}

// buildTicketPriceBreakdown groups tickets by class and calculates subtotals and total.
func (cs *SessionUsecase) HelperBuildPriceBreakdown(tickets []*entity.Ticket) ([]model.ClaimSessionTicketPricesResponse, float32) {
	ticketSummary := make(map[uint]*model.ClaimSessionTicketPricesResponse)
	var total float32

	for _, ticket := range tickets {
		classID := ticket.ClassID
		class := ticket.Class
		price := ticket.Price

		if _, exists := ticketSummary[classID]; !exists {
			ticketSummary[classID] = &model.ClaimSessionTicketPricesResponse{
				Class:    *mapper.TicketClassToSessionClassMapper.ToModel(&class),
				Price:    price,
				Quantity: 0,
				Subtotal: 0,
			}
		}

		ticketSummary[classID].Quantity++
		ticketSummary[classID].Subtotal += price
		total += price
	}

	summaryList := make([]model.ClaimSessionTicketPricesResponse, 0, len(ticketSummary))
	for _, entry := range ticketSummary {
		summaryList = append(summaryList, *entry)
	}

	return summaryList, total
}

// buildSessionListResponse maps a list of Session entities to response models.
func (cs *SessionUsecase) HelperBuildSessionListResponse(ctx context.Context, sessions []*entity.ClaimSession) []*model.ReadClaimSessionResponse {
	result := make([]*model.ReadClaimSessionResponse, len(sessions))
	for i, session := range sessions {
		err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
			tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
			if err != nil {
				return nil
			}
			result[i] = cs.HelperBuildResponse(session, tickets)
			return nil
		})

		if err != nil {
			return nil
		}
	}
	return result
}

func HelperExtractPassengerData(request *model.ClaimedSessionFillPassengerDataRequest) ([]uint, map[uint]model.ClaimedSessionPassengerDataInput) {
	ticketIDs := make([]uint, len(request.PassengerData))
	passengerMap := make(map[uint]model.ClaimedSessionPassengerDataInput)
	for i, data := range request.PassengerData {
		ticketIDs[i] = data.TicketID
		passengerMap[data.TicketID] = data
	}
	return ticketIDs, passengerMap
}

func (cs *SessionUsecase) HelperValidateAndUpdateTickets(
	tickets []*entity.Ticket,
	dataMap map[uint]model.ClaimedSessionPassengerDataInput,
	now time.Time,
) ([]uint, []model.ClaimedSessionTicketUpdateFailure, []*entity.Ticket) {

	retrievedTicketsMap := make(map[uint]*entity.Ticket)
	for _, ticket := range tickets {
		retrievedTicketsMap[ticket.ID] = ticket
	}

	var updatedIDs []uint
	var failed []model.ClaimedSessionTicketUpdateFailure
	var toUpdate []*entity.Ticket

	for id, data := range dataMap {
		ticket, exists := retrievedTicketsMap[id]
		if !exists {

			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Ticket not found in session"})
			continue
		}
		if ticket.Status != "pending_data_entry" {
			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: fmt.Sprintf("Status is %s", ticket.Status)})
			continue
		}
		if data.PassengerName == "" {
			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Passenger name required"})
			continue
		}
		if data.IDType == "" {
			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "ID Type required"})
			continue
		}
		if data.IDNumber == "" {
			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "ID Number required"})
			continue
		}
		if data.PassengerAge == 0 {
			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Passenger age required"})
			continue
		}
		if data.Address == "" {
			failed = append(failed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Passenger address required"})
			continue
		}
		ticket.PassengerName = &data.PassengerName
		ticket.PassengerAge = &data.PassengerAge
		ticket.Address = &data.Address
		ticket.IDType = &data.IDType
		ticket.IDNumber = &data.IDNumber
		ticket.SeatNumber = data.SeatNumber
		ticket.Status = "pending_payment"
		ticket.EntriesAt = &now
		toUpdate = append(toUpdate, ticket)
		updatedIDs = append(updatedIDs, ticket.ID)
	}

	return updatedIDs, failed, toUpdate
}
