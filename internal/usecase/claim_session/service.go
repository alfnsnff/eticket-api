package claim_session

import (
	"context"
	"errors"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClaimSessionUsecase struct {
	DB                   *gorm.DB
	SessionRepository    *repository.SessionRepository
	TicketRepository     *repository.TicketRepository
	ScheduleRepository   *repository.ScheduleRepository
	AllocationRepository *repository.AllocationRepository // Your AllocationRepository implements this
	ManifestRepository   *repository.ManifestRepository
	FareRepository       *repository.FareRepository
	BookingRepository    *repository.BookingRepository
}

func NewClaimSessionUsecase(
	db *gorm.DB,
	session_repository *repository.SessionRepository,
	ticket_repository *repository.TicketRepository,
	schedule_repository *repository.ScheduleRepository,
	allocation_repository *repository.AllocationRepository,
	manifest_repository *repository.ManifestRepository,
	fare_repository *repository.FareRepository,
	booking_repository *repository.BookingRepository,
) *ClaimSessionUsecase {
	return &ClaimSessionUsecase{
		DB:                   db,
		SessionRepository:    session_repository,
		TicketRepository:     ticket_repository,
		ScheduleRepository:   schedule_repository,
		AllocationRepository: allocation_repository,
		ManifestRepository:   manifest_repository,
		FareRepository:       fare_repository,
		BookingRepository:    booking_repository,
	}
}

func (cs *ClaimSessionUsecase) CreateClaimSession(ctx context.Context, request *model.ClaimedSessionLockTicketsRequest) (*model.ClaimedSessionLockTicketsResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	if request.ScheduleID == 0 {
		return nil, fmt.Errorf("mssing schedule ID")
	}
	if len(request.Items) == 0 {
		return nil, fmt.Errorf("missing request items")
	}
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

		count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(tx, request.ScheduleID, item.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"})
		if err != nil {
			return nil, fmt.Errorf("failed to count tickets: %w", err)
		}

		available := int64(cap.Quota) - count
		if available < int64(item.Quantity) {
			return nil, fmt.Errorf("not enough slots for class %d (Available: %d, Requested: %d)", item.ClassID, available, item.Quantity)
		}
	}

	var ticketsToBuild []*entity.Ticket

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

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session UUID: %w", err)
	}

	claimSession := &entity.ClaimSession{
		SessionID:  uuid.String(),
		ScheduleID: request.ScheduleID,
		ClaimedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(15 * time.Minute),
	}

	if err := cs.SessionRepository.Create(tx, claimSession); err != nil {
		return nil, fmt.Errorf("failed to create claim session: %w", err)
	}

	var claimedTicketIDs []uint
	for _, ticket := range ticketsToBuild {
		ticket.ClaimSessionID = &claimSession.ID
		claimedTicketIDs = append(claimedTicketIDs, ticket.ID)
	}

	if err := cs.TicketRepository.CreateBulk(tx, ticketsToBuild); err != nil {
		return nil, fmt.Errorf("failed to create tickets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.ClaimedSessionLockTicketsResponse{
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

	total, err := cs.SessionRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count claim sessions: %w", err)
	}

	claimsessions, err := cs.SessionRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all claim sessions: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return cs.HelperBuildSessionListResponse(ctx, claimsessions), int(total), nil
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

	session, err := cs.SessionRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return cs.HelperBuildResponse(session, tickets), nil
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

	session, err := cs.SessionRepository.GetByUUID(tx, sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets for session %d: %w", session.ID, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return cs.HelperBuildResponse(session, tickets), nil
}

func (cs *ClaimSessionUsecase) UpdateClaimSession(ctx context.Context, request *model.ClaimedSessionFillPassengerDataRequest, sessionID string) (*model.ClaimedSessionFillPassengerDataResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	if request.CustomerName == "" {
		return nil, errors.New("missing customer name")
	}

	if request.IDType == "" {
		return nil, errors.New("missing ID type")
	}

	if request.IDNumber == "" {
		return nil, errors.New("missing ID number")
	}

	if request.PhoneNumber == "" {
		return nil, errors.New("missing phone number")
	}

	if request.Email == "" {
		return nil, errors.New("missing email")
	}

	if len(request.TicketData) == 0 {
		return nil, errors.New("missing ticket data ")
	}

	// Build passenger map
	datas := make(map[uint]model.ClaimedSessionTicketDataInput)
	for _, data := range request.TicketData {
		datas[data.TicketID] = data
	}

	session, err := cs.SessionRepository.GetByUUIDWithLock(tx, sessionID, true)
	if err != nil {
		return nil, fmt.Errorf("get claim session failed: %w", err)
	}

	if session == nil || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("claim session not found or expired")
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
		return nil, fmt.Errorf("create booking failed: %w", err)
	}

	ticketsIds := make(map[uint]*entity.Ticket)
	for _, t := range tickets {
		tickets[t.ID] = t
	}

	var ticketToUpdate []*entity.Ticket
	var updatedTicketIDs []uint
	var total float32
	for id, data := range datas {
		ticket, ok := ticketsIds[id]
		if !ok {
			return nil, fmt.Errorf("ticket %d not found in session", id)
		}
		if ticket.Status != "pending_data_entry" {
			return nil, fmt.Errorf("ticket %d has invalid status: %s", id, ticket.Status)
		}
		if data.PassengerName == "" || data.PassengerAge == 0 || data.Address == "" {
			return nil, fmt.Errorf("missing passenger data for ticket %d", id)
		}

		ticket.PassengerName = &data.PassengerName
		ticket.PassengerAge = &data.PassengerAge
		ticket.Address = &data.Address
		ticket.BookingID = &booking.ID
		ticket.Status = "pending_payment"
		now := time.Now()
		ticket.EntriesAt = &now

		switch ticket.Type {
		case "passenger":
			if data.IDType == "" || data.IDNumber == "" {
				return nil, fmt.Errorf("missing ID info for passenger ticket %d", id)
			}
			ticket.IDType = &data.IDType
			ticket.IDNumber = &data.IDNumber

			count, err := cs.TicketRepository.CountByScheduleClassAndStatuses(
				tx, ticket.ScheduleID, ticket.ClassID, []string{"pending_data_entry", "pending_payment", "confirmed"},
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

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.ClaimedSessionFillPassengerDataResponse{
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

	session, err := cs.SessionRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return errors.New("session not found")
	}

	if err := cs.SessionRepository.Delete(tx, session); err != nil {
		return fmt.Errorf("failed to delete fare: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (cs *ClaimSessionUsecase) HelperBuildResponse(session *entity.ClaimSession, tickets []*entity.Ticket) *model.ReadClaimSessionResponse {
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

func (cs *ClaimSessionUsecase) HelperBuildTicketBreakdown(tickets []*entity.Ticket) []model.ClaimedSessionTicketDetailResponse {
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

func (cs *ClaimSessionUsecase) HelperBuildPriceBreakdown(tickets []*entity.Ticket) ([]model.ClaimSessionTicketPricesResponse, float32) {
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

func (cs *ClaimSessionUsecase) HelperBuildSessionListResponse(ctx context.Context, sessions []*entity.ClaimSession) []*model.ReadClaimSessionResponse {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	result := make([]*model.ReadClaimSessionResponse, len(sessions))
	for i, session := range sessions {
		tickets, err := cs.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			result[i] = cs.HelperBuildResponse(session, []*entity.Ticket{})                           // build with empty tickets
			fmt.Printf("Error processing session %s for list response: %v\n", session.SessionID, err) // Example logging
			continue                                                                                  // Don't let this specific error fail the whole list building
		}
		result[i] = cs.HelperBuildResponse(session, tickets)
	}
	finalResult := []*model.ReadClaimSessionResponse{}
	for _, res := range result {
		if res != nil {
			finalResult = append(finalResult, res)
		}
	}
	return finalResult
}
