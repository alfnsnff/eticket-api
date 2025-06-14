package claim_session

import (
	"context"
	"errors"
	"eticket-api/internal/client"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"eticket-api/pkg/payment"
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
	BookingRepository    *repository.BookingRepository
	TripayClient         *client.TripayClient
}

func NewSessionUsecase(
	tx *tx.TxManager,
	session_repository *repository.SessionRepository,
	ticket_repository *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	allocation_repository *repository.AllocationRepository,
	manifest_repository *repository.ManifestRepository,
	fare_repository *repository.FareRepository,
	booking_repository *repository.BookingRepository,
	tripay_client *client.TripayClient,
) *SessionUsecase {
	return &SessionUsecase{
		Tx:                   tx,
		SessionRepository:    session_repository,
		TicketRepository:     ticket_repository,
		ScheduleRepository:   schedule_repository,
		AllocationRepository: allocation_repository,
		ManifestRepository:   manifest_repository,
		FareRepository:       fare_repository,
		BookingRepository:    booking_repository,
		TripayClient:         tripay_client,
	}
}

func (cs *SessionUsecase) CreateSession(ctx context.Context, request *model.WriteClaimSessionRequest) error {
	session := mapper.SessionMapper.FromWrite(request)

	return cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return cs.SessionRepository.Create(tx, session)
	})
}

func (cs *SessionUsecase) GetAllSessions(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadClaimSessionResponse, int, error) {
	sessions := []*entity.ClaimSession{}
	var total int64
	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = cs.SessionRepository.Count(tx)
		if err != nil {
			return err
		}
		sessions, err = cs.SessionRepository.GetAll(tx, limit, offset, sort, search)
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
	if request.ScheduleID == 0 || len(request.Items) == 0 {
		return nil, fmt.Errorf("invalid claim request")
	}
	for _, item := range request.Items {
		if item.Quantity == 0 || item.ClassID == 0 {
			return nil, fmt.Errorf("missing request item field")
		}
	}

	var (
		claimedTicketIDs   []uint
		expiryTime         time.Time
		createdSessionUUID string
	)

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		schedule, err := cs.ScheduleRepository.GetByID(tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("failed to retrieve schedule: %w", err)
		}
		if schedule == nil {
			return fmt.Errorf("schedule not found")
		}

		// Validate availability per item
		for _, item := range request.Items {
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
		}

		now := time.Now()
		expiryTime = now.Add(15 * time.Minute)

		uuidVal, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("failed to generate session UUID: %w", err)
		}
		createdSessionUUID = uuidVal.String()

		var ticketsToBuild []*entity.Ticket

		for _, item := range request.Items {
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
					ClaimSessionID: nil, // Linked after session is created
				})
			}
		}

		newClaimSession := &entity.ClaimSession{
			SessionID:  createdSessionUUID,
			ScheduleID: request.ScheduleID,
			ClaimedAt:  now,
			ExpiresAt:  expiryTime,
		}

		if err := cs.SessionRepository.Create(tx, newClaimSession); err != nil {
			return fmt.Errorf("failed to create claim session: %w", err)
		}

		for _, ticket := range ticketsToBuild {
			ticket.ClaimSessionID = &newClaimSession.ID
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
	// Validate request
	if request.CustomerName == "" || request.IDType == "" || request.IDNumber == "" || request.PhoneNumber == "" || request.Email == "" {
		return nil, errors.New("invalid request: all customer fields are required")
	}
	if len(request.TicketData) == 0 {
		return nil, errors.New("invalid request: passenger data is required")
	}

	// Build passenger map
	passengerMap := make(map[uint]model.ClaimedSessionTicketDataInput)
	for _, data := range request.TicketData {
		passengerMap[data.TicketID] = data
	}

	var (
		finalUpdatedIDs []uint
		bookingID       uint
		xenditResp      payment.XenditResponse
		tripayResp      model.ReadTransactionResponse
		total           float32
	)

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		session, err := cs.SessionRepository.GetByUUIDWithLock(tx, sessionID, true)
		if err != nil {
			return fmt.Errorf("get claim session failed: %w", err)
		}
		if session == nil || session.ExpiresAt.Before(time.Now()) {
			return errors.New("claim session not found or expired")
		}

		tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("retrieve tickets failed: %w", err)
		}

		schedule, err := cs.ScheduleRepository.GetByID(tx, session.ScheduleID)
		if err != nil {
			return fmt.Errorf("retrieve schedule failed: %w", err)
		}

		orderID := utils.GenerateOrderID(
			fmt.Sprintf("%s-%s", *schedule.Route.DepartureHarbor.HarborAlias, *schedule.Route.ArrivalHarbor.HarborAlias),
			*schedule.Ship.ShipAlias, time.Now(),
		)

		booking := &entity.Booking{
			OrderID:      orderID,
			ScheduleID:   session.ScheduleID,
			IDType:       request.IDType,
			IDNumber:     request.IDNumber,
			PhoneNumber:  request.PhoneNumber,
			CustomerName: request.CustomerName,
			Email:        request.Email,
			Status:       "pending_payment",
		}
		if err := cs.BookingRepository.Create(tx, booking); err != nil {
			return fmt.Errorf("create booking failed: %w", err)
		}
		bookingID = booking.ID

		ticketsByID := make(map[uint]*entity.Ticket)
		for _, t := range tickets {
			ticketsByID[t.ID] = t
		}

		var updatedTickets []*entity.Ticket
		now := time.Now()

		for id, data := range passengerMap {
			ticket, ok := ticketsByID[id]
			if !ok {
				return fmt.Errorf("ticket %d not found in session", id)
			}
			if ticket.Status != "pending_data_entry" {
				return fmt.Errorf("ticket %d has invalid status: %s", id, ticket.Status)
			}
			if data.PassengerName == "" || data.PassengerAge == 0 || data.Address == "" {
				return fmt.Errorf("missing passenger data for ticket %d", id)
			}

			ticket.PassengerName = &data.PassengerName
			ticket.PassengerAge = &data.PassengerAge
			ticket.Address = &data.Address
			ticket.BookingID = &booking.ID
			ticket.Status = "pending_payment"
			ticket.EntriesAt = &now

			switch ticket.Type {
			case "passenger":
				if data.IDType == "" || data.IDNumber == "" {
					return fmt.Errorf("missing ID info for passenger ticket %d", id)
				}
				ticket.IDType = &data.IDType
				ticket.IDNumber = &data.IDNumber

				count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(
					tx, ticket.ScheduleID, ticket.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"},
				)
				if err != nil {
					return fmt.Errorf("seat number generation failed: %w", err)
				}
				seat := fmt.Sprintf("%s%d", *ticket.Class.ClassAlias, count+1)
				ticket.SeatNumber = &seat

			case "vehicle":
				if data.LicensePlate == nil || *data.LicensePlate == "" {
					return fmt.Errorf("missing license plate for vehicle ticket %d", id)
				}
				ticket.LicensePlate = data.LicensePlate
				ticket.SeatNumber = nil

			default:
				return fmt.Errorf("unsupported ticket type for ticket %d", id)
			}

			total += ticket.Price
			updatedTickets = append(updatedTickets, ticket)
			finalUpdatedIDs = append(finalUpdatedIDs, ticket.ID)
		}

		if len(updatedTickets) == len(tickets) {
			if err := cs.TicketRepository.UpdateBulk(tx, updatedTickets); err != nil {
				return fmt.Errorf("update tickets failed: %w", err)
			}

			// tripayResp, err = payment.CreateTripayPayment("QRIS", int(total), booking.CustomerName, booking.Email, *booking.OrderID, updatedTickets)
			// if err != nil {
			// 	return fmt.Errorf("create Tripay payment failed: %w", err)
			// }

			tripayResp, err = cs.TripayClient.CreatePayment("QRIS", int(total), booking.CustomerName, booking.Email, booking.PhoneNumber, booking.OrderID, updatedTickets)
			if err != nil {
				return fmt.Errorf("create Tripay payment failed: %w", err)
			}

			xenditResp, err = payment.CreateXenditPayment(int(total), booking.CustomerName, booking.Email, booking.OrderID, updatedTickets)
			if err != nil {
				return fmt.Errorf("create Xendit payment failed: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fill passenger data failed: %w", err)
	}

	return &model.ClaimedSessionFillPassengerDataResponse{
		BookingID:        bookingID,
		UpdatedTicketIDs: finalUpdatedIDs,
		Tripay:           tripayResp,
		Xendit:           xenditResp,
	}, nil
}

func (cs *SessionUsecase) DataEntry(ctx context.Context, request *model.ClaimedSessionFillPassengerDataRequest, sessionID string) (*model.ClaimedSessionFillPassengerDataResponse, error) {
	// Validate request
	if request.CustomerName == "" || request.IDType == "" || request.IDNumber == "" || request.PhoneNumber == "" || request.Email == "" {
		return nil, errors.New("invalid request: all customer fields are required")
	}
	if len(request.TicketData) == 0 {
		return nil, errors.New("invalid request: passenger data is required")
	}

	// Build passenger map
	passengerMap := make(map[uint]model.ClaimedSessionTicketDataInput)
	for _, data := range request.TicketData {
		passengerMap[data.TicketID] = data
	}

	var (
		finalUpdatedIDs []uint
		bookingID       uint
		total           float32
	)

	err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
		session, err := cs.SessionRepository.GetByUUIDWithLock(tx, sessionID, true)
		if err != nil {
			return fmt.Errorf("get claim session failed: %w", err)
		}
		if session == nil || session.ExpiresAt.Before(time.Now()) {
			return errors.New("claim session not found or expired")
		}

		tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("retrieve tickets failed: %w", err)
		}

		schedule, err := cs.ScheduleRepository.GetByID(tx, session.ScheduleID)
		if err != nil {
			return fmt.Errorf("retrieve schedule failed: %w", err)
		}

		orderID := utils.GenerateOrderID(
			fmt.Sprintf("%s-%s", *schedule.Route.DepartureHarbor.HarborAlias, *schedule.Route.ArrivalHarbor.HarborAlias),
			*schedule.Ship.ShipAlias, time.Now(),
		)

		booking := &entity.Booking{
			OrderID:      orderID,
			ScheduleID:   session.ScheduleID,
			IDType:       request.IDType,
			IDNumber:     request.IDNumber,
			PhoneNumber:  request.PhoneNumber,
			CustomerName: request.CustomerName,
			Email:        request.Email,
			Status:       "pending_payment",
		}
		if err := cs.BookingRepository.Create(tx, booking); err != nil {
			return fmt.Errorf("create booking failed: %w", err)
		}
		bookingID = booking.ID

		ticketsByID := make(map[uint]*entity.Ticket)
		for _, t := range tickets {
			ticketsByID[t.ID] = t
		}

		var updatedTickets []*entity.Ticket
		now := time.Now()

		for id, data := range passengerMap {
			ticket, ok := ticketsByID[id]
			if !ok {
				return fmt.Errorf("ticket %d not found in session", id)
			}
			if ticket.Status != "pending_data_entry" {
				return fmt.Errorf("ticket %d has invalid status: %s", id, ticket.Status)
			}
			if data.PassengerName == "" || data.PassengerAge == 0 || data.Address == "" {
				return fmt.Errorf("missing passenger data for ticket %d", id)
			}

			ticket.PassengerName = &data.PassengerName
			ticket.PassengerAge = &data.PassengerAge
			ticket.Address = &data.Address
			ticket.BookingID = &booking.ID
			ticket.Status = "pending_payment"
			ticket.EntriesAt = &now

			switch ticket.Type {
			case "passenger":
				if data.IDType == "" || data.IDNumber == "" {
					return fmt.Errorf("missing ID info for passenger ticket %d", id)
				}
				ticket.IDType = &data.IDType
				ticket.IDNumber = &data.IDNumber

				count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(
					tx, ticket.ScheduleID, ticket.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"},
				)
				if err != nil {
					return fmt.Errorf("seat number generation failed: %w", err)
				}
				seat := fmt.Sprintf("%s%d", *ticket.Class.ClassAlias, count+1)
				ticket.SeatNumber = &seat

			case "vehicle":
				if data.LicensePlate == nil || *data.LicensePlate == "" {
					return fmt.Errorf("missing license plate for vehicle ticket %d", id)
				}
				ticket.LicensePlate = data.LicensePlate
				ticket.SeatNumber = nil

			default:
				return fmt.Errorf("unsupported ticket type for ticket %d", id)
			}

			total += ticket.Price
			updatedTickets = append(updatedTickets, ticket)
			finalUpdatedIDs = append(finalUpdatedIDs, ticket.ID)
		}

		if len(updatedTickets) == len(tickets) {
			if err := cs.TicketRepository.UpdateBulk(tx, updatedTickets); err != nil {
				return fmt.Errorf("update tickets failed: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fill passenger data failed: %w", err)
	}

	return &model.ClaimedSessionFillPassengerDataResponse{
		BookingID:        bookingID,
		UpdatedTicketIDs: finalUpdatedIDs,
	}, nil
}

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

func (cs *SessionUsecase) HelperBuildTicketBreakdown(tickets []*entity.Ticket) []model.ClaimedSessionTicketDetailResponse {
	result := make([]model.ClaimedSessionTicketDetailResponse, len(tickets))
	for i, v := range tickets {
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

func (cs *SessionUsecase) HelperBuildPriceBreakdown(tickets []*entity.Ticket) ([]model.ClaimSessionTicketPricesResponse, float32) {
	ticketSummary := make(map[uint]*model.ClaimSessionTicketPricesResponse)
	var total float32

	for _, ticket := range tickets {
		classID := ticket.ClassID
		price := ticket.Price

		if _, exists := ticketSummary[classID]; !exists {
			var classModel model.ClaimSessionTicketClassItem
			if ticket.Class.ID != 0 { // Basic check
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

func (cs *SessionUsecase) HelperBuildSessionListResponse(ctx context.Context, sessions []*entity.ClaimSession) []*model.ReadClaimSessionResponse {
	result := make([]*model.ReadClaimSessionResponse, len(sessions))
	for i, session := range sessions {
		err := cs.Tx.Execute(ctx, func(tx *gorm.DB) error {
			tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
			if err != nil {
				result[i] = cs.HelperBuildResponse(session, []*entity.Ticket{}) // build with empty tickets
				return nil                                                      // Don't let this specific error fail the whole list building
			}
			result[i] = cs.HelperBuildResponse(session, tickets)
			return nil
		})

		if err != nil {
			fmt.Printf("Error processing session %s for list response: %v\n", session.SessionID, err) // Example logging

		}
	}
	finalResult := []*model.ReadClaimSessionResponse{}
	for _, res := range result {
		if res != nil {
			finalResult = append(finalResult, res)
		}
	}
	return finalResult
}
