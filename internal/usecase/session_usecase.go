package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionUsecase struct {
	DB                   *gorm.DB
	SessionRepository    *repository.SessionRepository
	TicketRepository     *repository.TicketRepository
	ScheduleRepository   *repository.ScheduleRepository
	AllocationRepository *repository.AllocationRepository // Your AllocationRepository implements this
	ManifestRepository   *repository.ManifestRepository
	FareRepository       *repository.FareRepository
}

func NewSessionUsecase(db *gorm.DB,
	sessionRepo *repository.SessionRepository,
	ticketRepo *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	allocation_repository *repository.AllocationRepository,
	manifest_repository *repository.ManifestRepository,
	fare_repository *repository.FareRepository,
) *SessionUsecase {
	return &SessionUsecase{
		DB:                   db,
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

	return tx.Execute(ctx, cs.DB, func(tx *gorm.DB) error {
		return cs.SessionRepository.Create(tx, session)
	})
}

func (cs *SessionUsecase) GetAllSessions(ctx context.Context) ([]*model.ReadClaimSessionResponse, error) {
	sessions := []*entity.ClaimSession{}

	err := tx.Execute(ctx, cs.DB, func(tx *gorm.DB) error {
		var err error
		sessions, err = cs.SessionRepository.GetAll(tx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all sessions: %w", err)
	}

	return cs.buildSessionListResponse(sessions), nil
}

func (cs *SessionUsecase) GetSessionByID(ctx context.Context, id uint) (*model.ReadClaimSessionResponse, error) {
	var session *entity.ClaimSession

	err := tx.Execute(ctx, cs.DB, func(tx *gorm.DB) error {
		var err error
		session, err = cs.SessionRepository.GetByID(tx, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	tickets, err := cs.TicketRepository.FindManyBySessionID(cs.DB, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
	}

	return cs.buildSessionResponse(session, tickets), nil
}

func (cs *SessionUsecase) UpdateSession(ctx context.Context, id uint, request *model.UpdateClaimSessionRequest) error {
	session := mapper.SessionMapper.FromUpdate(request)
	session.ID = id

	if session.ID == 0 {
		return fmt.Errorf("session ID cannot be zero")
	}

	return tx.Execute(ctx, cs.DB, func(tx *gorm.DB) error {
		return cs.SessionRepository.Update(tx, session)
	})
}

func (cs *SessionUsecase) DeleteSession(ctx context.Context, id uint) error {
	return tx.Execute(ctx, cs.DB, func(tx *gorm.DB) error {
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

	session, err := cs.SessionRepository.GetByUUID(cs.DB, sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session %s: %w", sessionUUID, err)
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	tickets, err := cs.TicketRepository.FindManyBySessionID(cs.DB, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
	}

	return cs.buildSessionResponse(session, tickets), nil
}

func (cs *SessionUsecase) LockBookingTickets(ctx context.Context, request *model.LockTicketsRequest) (*model.LockTicketsResponse, error) {
	if err := cs.validateLockRequest(request); err != nil {
		return nil, err
	}

	var claimedTicketIDs []uint
	var expiryTime time.Time
	var createdSessionUUID string // To hold the generated UUID

	err := tx.Execute(ctx, cs.DB, func(tx *gorm.DB) error {
		_, err := cs.lockAndCheckAvailability(tx, request)
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

		tickets, err := cs.buildTickets(tx, request)
		if err != nil {
			return err
		}

		newClaimSession := &entity.ClaimSession{
			SessionID: createdSessionUUID,
			ClaimedAt: now,
			ExpiresAt: expiryTime,
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

	return &model.LockTicketsResponse{
		ClaimedTicketIDs: claimedTicketIDs,
		ExpiresAt:        expiryTime,
		SessionID:        createdSessionUUID,
	}, nil
}

func (cs *SessionUsecase) validateLockRequest(request *model.LockTicketsRequest) error {
	if request.ScheduleID == 0 || len(request.Items) == 0 {
		return errors.New("invalid claim request")
	}
	return nil
}

func (cs *SessionUsecase) lockAndCheckAvailability(tx *gorm.DB, request *model.LockTicketsRequest) (map[uint]int64, error) {
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

func (cs *SessionUsecase) buildTickets(tx *gorm.DB, request *model.LockTicketsRequest) ([]*entity.Ticket, error) {
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
func (cs *SessionUsecase) buildSessionResponse(session *entity.ClaimSession, tickets []*entity.Ticket) *model.ReadClaimSessionResponse {
	var ticketDetails []model.ClaimSessionTicketPricesResponse
	var total float32

	if len(tickets) > 0 {
		ticketDetails, total = cs.buildTicketPriceBreakdown(tickets)
	} else {
		ticketDetails = []model.ClaimSessionTicketPricesResponse{}
		total = 0
	}

	return &model.ReadClaimSessionResponse{
		ID:          session.ID,
		SessionID:   session.SessionID,
		ClaimedAt:   session.ClaimedAt,
		ExpiresAt:   session.ExpiresAt,
		Tickets:     ticketDetails,
		TotalAmount: total,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}
}

// buildTicketPriceBreakdown groups tickets by class and calculates subtotals and total.
func (cs *SessionUsecase) buildTicketPriceBreakdown(tickets []*entity.Ticket) ([]model.ClaimSessionTicketPricesResponse, float32) {
	ticketSummary := make(map[uint]*model.ClaimSessionTicketPricesResponse)
	var total float32

	for _, ticket := range tickets {
		classID := ticket.ClassID
		price := ticket.Price

		if _, exists := ticketSummary[classID]; !exists {
			ticketSummary[classID] = &model.ClaimSessionTicketPricesResponse{
				ClassID:  classID,
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
func (cs *SessionUsecase) buildSessionListResponse(sessions []*entity.ClaimSession) []*model.ReadClaimSessionResponse {
	result := make([]*model.ReadClaimSessionResponse, len(sessions))
	for i, session := range sessions {
		tickets, err := cs.TicketRepository.FindManyBySessionID(cs.DB, session.ID)
		if err != nil {
			return nil
		}
		result[i] = cs.buildSessionResponse(session, tickets)
	}
	return result
}
