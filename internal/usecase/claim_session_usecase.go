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
	TripayClient           domain.TripayClient
}

func NewClaimSessionUsecase(
	db *gorm.DB,
	claim_session_repository domain.ClaimSessionRepository,
	claim_item_repository domain.ClaimItemRepository, // Assuming you have a ClaimItemRepository
	ticket_repository domain.TicketRepository,
	schedule_repository domain.ScheduleRepository,
	booking_repository domain.BookingRepository,
	quota_repository domain.QuotaRepository,
	tripay_client domain.TripayClient,
) *ClaimSessionUsecase {
	return &ClaimSessionUsecase{
		DB:                     db,
		ClaimSessionRepository: claim_session_repository,
		ClaimItemRepository:    claim_item_repository, // Assuming you have a ClaimItemRepository
		TicketRepository:       ticket_repository,
		ScheduleRepository:     schedule_repository,
		BookingRepository:      booking_repository,
		QuotaRepository:        quota_repository,
		TripayClient:           tripay_client,
	}
}

func (cs *ClaimSessionUsecase) TESTCreateClaimSession(
	ctx context.Context,
	request *model.TESTWriteClaimSessionRequest,
) (*model.TESTReadClaimSessionLockResponse, error) {
	tx := cs.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Step 1: Validate schedule existence
	schedule, err := cs.ScheduleRepository.FindByID(tx, request.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("retrieve schedule: %w", err)
	}
	if schedule == nil {
		return nil, errs.ErrNotFound
	}

	// Step 2: Fetch and map quotas
	quotas, err := cs.QuotaRepository.FindByScheduleID(tx, request.ScheduleID)
	if err != nil {

		return nil, fmt.Errorf("fetch quotas: %w", err)
	}
	quotaByClass := make(map[uint]*domain.Quota, len(quotas))
	for i := range quotas {
		quotaByClass[quotas[i].ClassID] = quotas[i]
	}

	// Step 3: Fetch active claim sessions and accumulate usage
	sessions, err := cs.ClaimSessionRepository.FindByScheduleID(tx, request.ScheduleID)
	if err != nil {

		return nil, fmt.Errorf("load active sessions: %w", err)
	}
	usedByClass := make(map[uint]int64)
	for i := range sessions {
		for j := range sessions[i].ClaimItems {
			item := sessions[i].ClaimItems[j]
			usedByClass[item.ClassID] += int64(item.Quantity)
		}
	}

	// Step 4: Validate quota availability
	for i := range request.Items {
		classID := request.Items[i].ClassID
		quota, exists := quotaByClass[classID]
		if !exists {
			tx.Rollback()
			return nil, fmt.Errorf("quota not found for class %d", classID)
		}
		used := usedByClass[classID]
		requested := int64(request.Items[i].Quantity)
		if used+requested > int64(quota.Quota) {
			tx.Rollback()
			return nil, fmt.Errorf("quota exceeded for class %d", classID)
		}
	}

	claimItems := make([]domain.ClaimItem, len(request.Items))
	for i, item := range request.Items {
		claimItems[i] = domain.ClaimItem{
			ClassID:  item.ClassID,
			Quantity: item.Quantity,
		}
	}

	// Step 5: Create claim session
	claimSession := &domain.ClaimSession{
		SessionID:  uuid.NewString(),
		ScheduleID: request.ScheduleID,
		Status:     enum.ClaimSessionPendingData.String(),
		ExpiresAt:  time.Now().Add(16 * time.Minute),
		ClaimItems: claimItems, // attach here
	}
	if err := cs.ClaimSessionRepository.Insert(tx, claimSession); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create session: %w", err)
	}

	// Step 7: Commit
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// Step 8: Response
	return &model.TESTReadClaimSessionLockResponse{
		SessionID: claimSession.SessionID,
		ExpiresAt: claimSession.ExpiresAt,
	}, nil
}

func (cd *ClaimSessionUsecase) TESTUpdateClaimSession(
	ctx context.Context,
	request *model.TESTWriteClaimSessionDataEntryRequest,
	sessionID string,
) (string, error) {
	tx := cd.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Load session
	session, err := cd.ClaimSessionRepository.FindBySessionID(tx, sessionID)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("get claim session failed: %w", err)
	}
	if session == nil {
		tx.Rollback()
		return "", errs.ErrNotFound
	}
	if session.ExpiresAt.Before(time.Now()) {
		tx.Rollback()
		return "", errors.New("claim session expired")
	}
	if session.Status == "LOCKED" {
		tx.Rollback()
		return "", fmt.Errorf("claim session has invalid status: %s", session.Status)
	}

	// Generate order ID
	orderID := utils.GenerateOrderID(
		session.Schedule.DepartureHarbor.HarborAlias,
		session.Schedule.ArrivalHarbor.HarborAlias,
		session.Schedule.Ship.ShipAlias,
		time.Now(),
	)

	// Create booking
	booking := &domain.Booking{
		OrderID:      &orderID,
		ScheduleID:   session.ScheduleID,
		IDType:       request.IDType,
		IDNumber:     request.IDNumber,
		PhoneNumber:  request.PhoneNumber,
		CustomerName: request.CustomerName,
		Email:        request.Email,
		Status:       "PENDING_PAYMENT",
		ExpiresAt:    time.Now().Add(13 * time.Minute), // Set expiration for 13 minutes
	}
	if err := cd.BookingRepository.Insert(tx, booking); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("create booking failed: %w", err)
	}

	// Fetch quota and map by ClassID
	quotas, err := cd.QuotaRepository.FindByScheduleID(tx, session.ScheduleID)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("fetch quotas failed: %w", err)
	}
	quotaByClass := make(map[uint]*domain.Quota, len(quotas))
	for _, q := range quotas {
		quotaByClass[q.ClassID] = q
	}

	// Organize passenger data by ClassID
	dataQueue := make(map[uint][]model.TESTClaimSessionTicketDataEntry)
	for _, d := range request.TicketData {
		dataQueue[d.ClassID] = append(dataQueue[d.ClassID], d)
	}

	// Build ticket list
	var tickets []*domain.Ticket
	var amounts float64
	for _, item := range session.ClaimItems {
		quota, ok := quotaByClass[item.ClassID]
		if !ok {
			tx.Rollback()
			return "", fmt.Errorf("quota not found for class %d", item.ClassID)
		}

		classData := dataQueue[item.ClassID]
		if len(classData) < item.Quantity {
			tx.Rollback()
			return "", fmt.Errorf("not enough ticket data for class %d", item.ClassID)
		}

		for i := 0; i < item.Quantity; i++ {
			data := classData[i]
			tickets = append(tickets, &domain.Ticket{
				BookingID:       &booking.ID,
				ClassID:         item.ClassID,
				Price:           quota.Price,
				Type:            quota.Class.Type,
				PassengerName:   &data.PassengerName,
				IDType:          &data.IDType,
				IDNumber:        &data.IDNumber,
				PassengerAge:    &data.PassengerAge,
				PassengerGender: &data.PassengerGender,
				Address:         &data.Address,
				SeatNumber:      data.SeatNumber,
				LicensePlate:    data.LicensePlate,
				ScheduleID:      session.ScheduleID,
				ClaimSessionID:  &session.ID,
			})
			amounts += quota.Price
		}

		// Trim used data
		dataQueue[item.ClassID] = classData[item.Quantity:]
	}

	// Insert tickets
	if err := cd.TicketRepository.InsertBulk(tx, tickets); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create tickets: %w", err)
	}

	orderItems := make([]model.OrderItem, len(tickets))
	for i, ticket := range tickets {
		orderItems[i] = mapper.TicketToItem(ticket)
	}

	payload := &model.WriteTransactionRequest{
		Method:        request.PaymentMethod,
		Amount:        int(amounts), // Convert to integer cents
		CustomerName:  booking.CustomerName,
		CustomerEmail: booking.Email,
		CustomerPhone: booking.PhoneNumber,
		MerchantRef:   *booking.OrderID,
		OrderItems:    orderItems,
		CallbackUrl:   "https://example.com/callback",
		ReturnUrl:     "https://example.com/callback",
		ExpiredTime:   int(time.Now().Add(30 * time.Minute).Unix()),
	}

	payment, err := cd.TripayClient.CreatePayment(payload)
	if err != nil {
		return "", fmt.Errorf("create Tripay payment failed: %w", err)
	}

	booking.ReferenceNumber = &payment.Reference
	if err := cd.BookingRepository.Update(tx, booking); err != nil {
		return "", fmt.Errorf("failed to update booking with reference number: %w", err)
	}

	// Decrement quota usage
	for _, item := range session.ClaimItems {
		quota, ok := quotaByClass[item.ClassID]
		if !ok {
			tx.Rollback()
			return "", fmt.Errorf("quota not found for class %d", item.ClassID)
		}

		if quota.Quota < item.Quantity {
			tx.Rollback()
			return "", fmt.Errorf("not enough quota for class %d", item.ClassID)
		}

		quota.Quota -= item.Quantity

		if err := cd.QuotaRepository.Update(tx, quota); err != nil {
			tx.Rollback()
			return "", fmt.Errorf("failed to update quota usage: %w", err)
		}
	}

	// Mark session completed
	session.Status = "LOCKED"
	if err := cd.ClaimSessionRepository.Update(tx, session); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to update session: %w", err)
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("commit transaction failed: %w", err)
	}

	return *booking.OrderID, nil
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
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve schedule: %w", err)
	}
	if schedule == nil {
		tx.Rollback()
		return nil, errs.ErrNotFound
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate session UUID: %w", err)
	}

	claimSession := &domain.ClaimSession{
		SessionID:  uuid.String(),
		ScheduleID: request.ScheduleID,
		Status:     enum.ClaimSessionPendingData.String(),
		ExpiresAt:  time.Now().Add(13 * time.Minute),
	}

	// ⬇️ Insert the ClaimSession first so we get claimSession.ID populated
	if err := cs.ClaimSessionRepository.Insert(tx, claimSession); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create claim session: %w", err)
	}

	var ticketsToBuild []*domain.Ticket
	for _, item := range request.Items {
		if item.ClassID == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("class is zero value in request item")
		}
		quota, err := cs.QuotaRepository.FindByScheduleIDAndClassID(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to get quota capacity: %w", err)
		}
		if quota == nil {
			tx.Rollback()
			return nil, fmt.Errorf("quota not found for class %d", item.ClassID)
		}
		occupied, err := cs.TicketRepository.CountByScheduleIDAndClassIDWithStatus(tx, request.ScheduleID, item.ClassID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to count tickets: %w", err)
		}

		available := int64(quota.Quota) - occupied
		if available < int64(item.Quantity) {
			tx.Rollback()
			return nil, fmt.Errorf("not enough slots for class %d (Available: %d, Requested: %d)", item.ClassID, available, item.Quantity)
		}

		for i := 0; i < int(item.Quantity); i++ {
			ticketsToBuild = append(ticketsToBuild, &domain.Ticket{
				ScheduleID:     request.ScheduleID,
				ClassID:        item.ClassID,
				Price:          quota.Price,
				Type:           quota.Class.Type,
				ClaimSessionID: &claimSession.ID,
				IsCheckedIn:    false,
			})
		}
	}

	if err := cs.TicketRepository.InsertBulk(tx, ticketsToBuild); err != nil {
		tx.Rollback()
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

func (cs *ClaimSessionUsecase) GetClaimSessionByID(ctx context.Context, id uint) (*model.TESTReadClaimSessionResponse, error) {
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

	return mapper.TESTClaimSessionToResponse(session), nil
}

func (cs *ClaimSessionUsecase) GetBySessionID(ctx context.Context, sessionUUID string) (*model.TESTReadClaimSessionResponse, error) {
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
	fmt.Printf("Fetched %d claim items", len(session.ClaimItems))

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.TESTClaimSessionToResponse(session), nil
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
		case "Passenger":
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

		case "Vehicle":
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
