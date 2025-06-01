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

// --- End of Example Interface Definitions ---

func (cs *SessionUsecase) SessionLockTickets(ctx context.Context, request *model.ClaimedSessionLockTicketsRequest) (*model.ClaimedSessionLockTicketsResponse, error) {
	// Inlined HelperValidateLockRequest
	if request.ScheduleID == 0 || len(request.Items) == 0 {
		return nil, fmt.Errorf("invalid claim request")
	}
	for _, item := range request.Items {
		if item.Quantity == 0 || item.Type == "" || item.ClassID == 0 {
			return nil, fmt.Errorf("missing request item field")
		}
	}

	var claimedTicketIDs []uint
	var expiryTime time.Time
	var createdSessionUUID string // To hold the generated UUID

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {

		sch, err := cs.ScheduleRepository.GetByID(tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("failed to lock capacity: %w", err)
		}
		if sch == nil {
			return fmt.Errorf("schedule not found")
		}
		// Inlined HelperLockAndCheckAvailability
		// The 'checks' map was not used further in the original SessionLockTickets, so its direct assignment is omitted.
		// The logic for checking availability remains.
		for _, item := range request.Items {
			// item validation was already done by the inlined HelperValidateLockRequest,
			// but keeping it here as it was in HelperLockAndCheckAvailability for logical grouping if this block was a separate func.
			// For direct inlining, it's somewhat redundant but harmless.
			if item.Quantity == 0 || item.Type == "" || item.ClassID == 0 {
				return fmt.Errorf("failed to lock capacity: missing request item field")
			}

			cap, err := cs.AllocationRepository.LockByScheduleAndClass(tx, request.ScheduleID, item.ClassID)
			if err != nil {
				return fmt.Errorf("failed to lock capacity: %w", err)
			}
			if cap == nil {
				return fmt.Errorf("allocation not found for class %d schedule %d", item.ClassID, request.ScheduleID)
			}

			count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(tx, request.ScheduleID, item.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"})
			if err != nil {
				return fmt.Errorf("failed to count tickets: %w", err)
			}

			available := int64(cap.Quota) - count
			if available < int64(item.Quantity) {
				return fmt.Errorf("not enough slots for class %d (Available: %d, Requested: %d)", item.ClassID, available, item.Quantity)
			}
			// checks[item.ClassID] = available // 'checks' map was not used
		}
		// End of Inlined HelperLockAndCheckAvailability logic

		now := time.Now()
		expiryTime = time.Now().Add(15 * time.Minute)

		sessionUUID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("failed to generate session UUID: %w", err) // This will trigger rollback
		}
		createdSessionUUID = sessionUUID.String() // Store for response

		// Inlined HelperBuildTickets
		schedule, err := cs.ScheduleRepository.GetByID(tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("failed to get schedule: %w", err)
		}

		var ticketsToBuild []*entity.Ticket
		for _, item := range request.Items {
			if item.Quantity == 0 {
				continue
			}

			manifest, err := cs.ManifestRepository.GetByShipAndClass(tx, schedule.ShipID, item.ClassID)
			if err != nil || manifest == nil {
				return fmt.Errorf("manifest missing for ship %d, class %d", schedule.ShipID, item.ClassID)
			}

			fare, err := cs.FareRepository.GetByManifestAndRoute(tx, manifest.ID, schedule.RouteID)
			if err != nil || fare == nil {
				return fmt.Errorf("fare missing for manifest %d, route %d", manifest.ID, schedule.RouteID)
			}

			for i := 0; i < int(item.Quantity); i++ {
				ticketsToBuild = append(ticketsToBuild, &entity.Ticket{
					ScheduleID:     request.ScheduleID,
					ClassID:        item.ClassID,
					Status:         "pending_data_entry",
					Price:          fare.TicketPrice,
					Type:           manifest.Class.Type,
					ClaimSessionID: nil, // Will be set after ClaimSession is created
				})
			}
		}
		// End of Inlined HelperBuildTickets logic
		// 'ticketsToBuild' is the result of the inlined HelperBuildTickets

		newClaimSession := &entity.ClaimSession{
			SessionID:  createdSessionUUID,
			ScheduleID: request.ScheduleID,
			ClaimedAt:  now,
			ExpiresAt:  expiryTime,
		}

		err = cs.SessionRepository.Create(tx, newClaimSession) // Use txDB
		if err != nil {
			return fmt.Errorf("failed to create claim session: %w", err) // This will trigger rollback
		}

		// Link the newly created tickets to the ClaimSession
		for _, ticket := range ticketsToBuild {
			ticket.ClaimSessionID = &newClaimSession.ID // Set the FK to the new ClaimSession ID
		}

		if err := cs.TicketRepository.CreateBulk(tx, ticketsToBuild); err != nil {
			return fmt.Errorf("failed to create tickets: %w", err)
		}

		for _, t := range ticketsToBuild {
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
	if len(request.TicketData) == 0 {
		return nil, errors.New("invalid request: passenger data are required")
	}

	// Inlined HelperExtractPassengerData
	// ticketIDsFromHelper was not used, so we only extract passengerMap
	passengerMap := make(map[uint]model.ClaimedSessionTicketDataInput)
	for _, data := range request.TicketData {
		passengerMap[data.TicketID] = data
	}
	// End of Inlined HelperExtractPassengerData logic

	var finalUpdatedIDs []uint
	var finalFailed []model.ClaimedSessionTicketUpdateFailure

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

		ticketsFromDB, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets: %w", err)
		}

		// Inlined HelperValidateAndUpdateTickets
		retrievedTicketsMap := make(map[uint]*entity.Ticket)
		for _, ticket := range ticketsFromDB {
			retrievedTicketsMap[ticket.ID] = ticket
		}

		var currentUpdatedIDs []uint
		var currentFailed []model.ClaimedSessionTicketUpdateFailure
		var ticketsToUpdate []*entity.Ticket

		for id, data := range passengerMap { // passengerMap from inlined HelperExtractPassengerData
			ticket, exists := retrievedTicketsMap[id]
			if !exists {
				currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Ticket not found in session"})
				continue
			}
			if ticket.Status != "pending_data_entry" {
				currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: fmt.Sprintf("Status is %s", ticket.Status)})
				continue
			}

			switch ticket.Type {
			case "passenger":
				if data.PassengerName == "" {
					currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Passenger name required"})
					continue
				}
				if data.IDType == "" {
					currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "ID Type required"})
					continue
				}
				if data.IDNumber == "" {
					currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "ID Number required"})
					continue
				}
				if data.PassengerAge == 0 {
					currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Passenger age required"})
					continue
				}
				if data.Address == "" {
					currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Passenger address required"})
					continue
				}
				ticket.PassengerName = &data.PassengerName
				ticket.PassengerAge = &data.PassengerAge
				ticket.Address = &data.Address
				ticket.IDType = &data.IDType
				ticket.IDNumber = &data.IDNumber

				count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(tx, ticket.ScheduleID, ticket.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"})
				if err != nil {
					return fmt.Errorf("failed to create seat number: %w", err)
				}
				seatNumStr := fmt.Sprintf("%s%d", ticket.Class.ClassName, count+1)
				ticket.SeatNumber = &seatNumStr

			case "vehicle":
				if data.LicensePlate == nil || *data.LicensePlate == "" {
					currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "License plate required for vehicle ticket"})
					continue
				}
				ticket.LicensePlate = data.LicensePlate
				ticket.SeatNumber = nil
			default:
				currentFailed = append(currentFailed, model.ClaimedSessionTicketUpdateFailure{TicketID: id, Reason: "Unsupported ticket type"})
				continue
			}

			ticket.Status = "pending_payment"
			ticket.EntriesAt = &now // now is from the outer scope (SessionDataEntry)
			ticketsToUpdate = append(ticketsToUpdate, ticket)
			currentUpdatedIDs = append(currentUpdatedIDs, ticket.ID)
		}
		// End of Inlined HelperValidateAndUpdateTickets logic

		finalUpdatedIDs = currentUpdatedIDs
		finalFailed = currentFailed

		if len(ticketsToUpdate) > 0 { // ticketsToUpdate is the result of inlined HelperValidateAndUpdateTickets
			err = cs.TicketRepository.UpdateBulk(tx, ticketsToUpdate)
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
		UpdatedTicketIDs: finalUpdatedIDs,
		FailedTickets:    finalFailed,
	}, nil
}

// --- Remaining Helper functions related to response building ---
// These helpers were not directly called by SessionLockTickets or SessionDataEntry in the provided snippet,
// or are used by other functions (like a potential GetSessionDetails or ListSessions).
// They are kept separate as their primary role is response shaping.

// buildSessionResponse generates a consistent ReadClaimSessionResponse, optionally using ticket info.
func (cs *SessionUsecase) HelperBuildResponse(session *entity.ClaimSession, tickets []*entity.Ticket) *model.ReadClaimSessionResponse {
	var ticketPrices []model.ClaimSessionTicketPricesResponse
	var ticketDetails []model.ClaimedSessionTicketDetailResponse
	var total float32

	if len(tickets) > 0 {
		ticketPrices, total = cs.HelperBuildPriceBreakdown(tickets)
		ticketDetails = cs.HelperBuildTicketBreakdown(tickets)
	} else {
		ticketPrices = []model.ClaimSessionTicketPricesResponse{}    // Ensure empty slice, not nil
		ticketDetails = []model.ClaimedSessionTicketDetailResponse{} // Ensure empty slice, not nil
		total = 0
	}

	// Assuming session.Schedule is preloaded or handled correctly
	// If session.Schedule might be zero/uninitialized, you'd need a check here
	var scheduleModel model.ClaimSessionSchedule // Default to zero value if session.Schedule is not populated
	if session.Schedule.ID != 0 {                // Basic check; adjust as needed based on your entity structure
		scheduleModel = *mapper.ScheduleSessionMapper.ToModel(&session.Schedule)
	}

	return &model.ReadClaimSessionResponse{
		ID:          session.ID,
		SessionID:   session.SessionID,
		ScheduleID:  session.ScheduleID,
		Schedule:    scheduleModel,
		ClaimedAt:   session.ClaimedAt,
		ExpiresAt:   session.ExpiresAt,
		Prices:      ticketPrices,
		Tickets:     ticketDetails,
		TotalAmount: total,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}
}

// HelperBuildTicketBreakdown groups tickets by class and calculates subtotals and total.
func (cs *SessionUsecase) HelperBuildTicketBreakdown(tickets []*entity.Ticket) []model.ClaimedSessionTicketDetailResponse {
	result := make([]model.ClaimedSessionTicketDetailResponse, len(tickets))
	for i, v := range tickets {
		// Assuming v.Class is preloaded or handled correctly
		// If v.Class might be zero/uninitialized, you'd need a check here
		var classModel model.ClaimSessionTicketClassItem // Default to zero value
		if v.Class.ID != 0 {                             // Basic check
			classModel = *mapper.TicketClassToSessionClassMapper.ToModel(&v.Class)
		}

		result[i] = model.ClaimedSessionTicketDetailResponse{
			TicketID: v.ID,
			Class:    classModel,
			Price:    v.Price,
			Type:     v.Type,
		}
	}
	return result
}

// HelperBuildPriceBreakdown groups tickets by class and calculates subtotals and total.
func (cs *SessionUsecase) HelperBuildPriceBreakdown(tickets []*entity.Ticket) ([]model.ClaimSessionTicketPricesResponse, float32) {
	ticketSummary := make(map[uint]*model.ClaimSessionTicketPricesResponse)
	var total float32

	for _, ticket := range tickets {
		classID := ticket.ClassID
		// class := ticket.Class // This is the entity.Class
		price := ticket.Price

		if _, exists := ticketSummary[classID]; !exists {
			// Assuming ticket.Class is preloaded or handled correctly
			var classModel model.ClaimSessionTicketClassItem // Default to zero value
			if ticket.Class.ID != 0 {                        // Basic check
				classModel = *mapper.TicketClassToSessionClassMapper.ToModel(&ticket.Class)
			}
			ticketSummary[classID] = &model.ClaimSessionTicketPricesResponse{
				Class:    classModel,
				Price:    price, // This is price per ticket
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

// HelperBuildSessionListResponse maps a list of Session entities to response models.
func (cs *SessionUsecase) HelperBuildSessionListResponse(ctx context.Context, sessions []*entity.ClaimSession) []*model.ReadClaimSessionResponse {
	result := make([]*model.ReadClaimSessionResponse, len(sessions))
	for i, session := range sessions {
		// This transaction wrapper here might be problematic if the outer function already runs in a transaction.
		// Or, if this is meant to be a separate unit of work for each session, it's fine.
		// Consider the transactional boundaries carefully. For simplicity of this refactor, I'm keeping it.
		err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
			// It's generally better to fetch tickets for ALL sessions in one go outside the loop
			// and then pass them to HelperBuildResponse, to avoid N+1 queries within Tx.Execute.
			// However, sticking to the original structure for this refactoring:
			tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
			if err != nil {
				// Log error or handle; returning nil might hide issues
				// For now, if tickets can't be fetched, an empty response for that session part is built.
				result[i] = cs.HelperBuildResponse(session, []*entity.Ticket{}) // build with empty tickets
				return nil                                                      // Don't let this specific error fail the whole list building
			}
			result[i] = cs.HelperBuildResponse(session, tickets)
			return nil
		})

		if err != nil {
			// If Tx.Execute fails for some reason (not FindManyBySessionID, which is handled inside)
			// Log this error. Depending on requirements, you might want to stop or continue.
			// For now, this would mean result[i] might be nil if the Execute itself failed before HelperBuildResponse was called.
			// To be safe, initialize result[i] to a minimal representation or skip.
			// Given the current structure, if err is not nil here, result[i] might not be set.
			// This part of the logic might need more robust error handling.
			// For example, if Tx.Execute fails, what should be in result[i]?
			// Perhaps:
			// result[i] = &model.ReadClaimSessionResponse{ SessionID: session.SessionID, Error: "Failed to load details" }
			// For now, it results in a nil entry in the list if Tx.Execute fails.
			// A simple `return nil` here would discard the entire list if one item fails.
			// Better to log and continue or handle partial results.
			fmt.Printf("Error processing session %s for list response: %v\n", session.SessionID, err) // Example logging
			// To ensure the slot isn't nil, you could do:
			// if result[i] == nil {
			//    result[i] = cs.HelperBuildResponse(session, []*entity.Ticket{}) // Default response
			// }
		}
	}
	// Filter out nil results if any Execute failed catastrophically for an item
	finalResult := []*model.ReadClaimSessionResponse{}
	for _, res := range result {
		if res != nil {
			finalResult = append(finalResult, res)
		}
	}
	return finalResult
}
