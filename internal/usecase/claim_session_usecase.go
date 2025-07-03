package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/client"
	enum "eticket-api/internal/common/enums"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"eticket-api/pkg/gotann"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ClaimSessionUsecase struct {
	Transactor             transact.Transactor
	ClaimSessionRepository domain.ClaimSessionRepository
	ClaimItemRepository    domain.ClaimItemRepository // Assuming you have a ClaimItemRepository
	TicketRepository       domain.TicketRepository
	ScheduleRepository     domain.ScheduleRepository
	BookingRepository      domain.BookingRepository
	QuotaRepository        domain.QuotaRepository
	TripayClient           domain.TripayClient
	Mailer                 mailer.Mailer // Assuming you have a Mailer interface for sending emails
}

func NewClaimSessionUsecase(
	transactor transact.Transactor,
	claim_session_repository domain.ClaimSessionRepository,
	claim_item_repository domain.ClaimItemRepository, // Assuming you have a ClaimItemRepository
	ticket_repository domain.TicketRepository,
	schedule_repository domain.ScheduleRepository,
	booking_repository domain.BookingRepository,
	quota_repository domain.QuotaRepository,
	tripay_client domain.TripayClient,
	mailer mailer.Mailer, // Assuming you have a Mailer interface for sending emails
) *ClaimSessionUsecase {
	return &ClaimSessionUsecase{
		Transactor:             transactor,
		ClaimSessionRepository: claim_session_repository,
		ClaimItemRepository:    claim_item_repository,
		TicketRepository:       ticket_repository,
		ScheduleRepository:     schedule_repository,
		BookingRepository:      booking_repository,
		QuotaRepository:        quota_repository,
		TripayClient:           tripay_client,
		Mailer:                 mailer, // Initialize the Mailer
	}
}

func (uc *ClaimSessionUsecase) LockClaimSession(
	ctx context.Context,
	request *model.TESTWriteClaimSessionRequest,
) (*model.TESTReadClaimSessionLockResponse, error) {

	var claimSession *domain.ClaimSession

	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		// Step 1: Validate schedule existence
		schedule, err := uc.ScheduleRepository.FindByID(ctx, tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("retrieve schedule: %w", err)
		}
		if schedule == nil {
			return errs.ErrNotFound
		}

		// Step 2: Fetch and map quotas
		quotas, err := uc.QuotaRepository.FindByScheduleID(ctx, tx, request.ScheduleID)
		if err != nil {

			return fmt.Errorf("fetch quotas: %w", err)
		}
		quotaByClass := make(map[uint]*domain.Quota, len(quotas))
		for i := range quotas {
			quotaByClass[quotas[i].ClassID] = quotas[i]
		}

		// Step 3: Fetch active claim sessions and accumulate usage
		sessions, err := uc.ClaimSessionRepository.FindActiveByScheduleID(ctx, tx, request.ScheduleID)
		if err != nil {
			return fmt.Errorf("load active sessions: %w", err)
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
				return fmt.Errorf("quota not found for class %d", classID)
			}
			used := usedByClass[classID]
			requested := int64(request.Items[i].Quantity)
			if used+requested > int64(quota.Quota) {
				return fmt.Errorf("quota exceeded for class %d", classID)
			}
		}

		claimItems := make([]domain.ClaimItem, len(request.Items))
		for i, item := range request.Items {
			quota, exists := quotaByClass[item.ClassID]
			if !exists {
				return fmt.Errorf("quota not found for class %d", item.ClassID)
			}
			subtotal := float64(item.Quantity) * quota.Price
			claimItems[i] = domain.ClaimItem{
				ClassID:  item.ClassID,
				Quantity: item.Quantity,
				Subtotal: subtotal, // <-- set subtotal by quota price
			}
		}

		// Step 5: Create claim session
		claimSession = &domain.ClaimSession{
			SessionID:  uuid.NewString(),
			ScheduleID: request.ScheduleID,
			Status:     enum.ClaimSessionPending.String(),
			ExpiresAt:  time.Now().Add(16 * time.Minute),
			ClaimItems: claimItems, // attach here
		}
		if err := uc.ClaimSessionRepository.Insert(ctx, tx, claimSession); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("create session: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("execute transaction: %w", err)
	}

	// Step 8: Response
	return &model.TESTReadClaimSessionLockResponse{
		SessionID: claimSession.SessionID,
		ExpiresAt: claimSession.ExpiresAt,
	}, nil
}

func (cd *ClaimSessionUsecase) EntryClaimSession(
	ctx context.Context,
	request *model.TESTWriteClaimSessionDataEntryRequest,
	sessionID string,
) (*model.TESTReadClaimSessionDataEntryResponse, error) {

	var booking *domain.Booking
	if err := cd.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		session, err := cd.ClaimSessionRepository.FindBySessionID(ctx, tx, sessionID)
		if err != nil {

			return fmt.Errorf("get claim session failed: %w", err)
		}
		if session == nil {

			return errs.ErrNotFound
		}
		if session.ExpiresAt.Before(time.Now()) {
			return errors.New("claim session expired")
		}
		// Generate order ID
		orderID := utils.GenerateOrderID(session.Schedule.DepartureHarbor.HarborAlias)

		// Create booking
		booking = &domain.Booking{
			OrderID:      orderID,
			ScheduleID:   session.ScheduleID,
			IDType:       request.IDType,
			IDNumber:     request.IDNumber,
			PhoneNumber:  request.PhoneNumber,
			CustomerName: request.CustomerName,
			Email:        request.Email,
			Status:       enum.ClaimSessionPending.String(),
			ExpiresAt:    time.Now().Add(13 * time.Minute), // Set expiration for 13 minutes
		}
		if err := cd.BookingRepository.Insert(ctx, tx, booking); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("create booking failed: %w", err)
		}

		// Fetch quota and map by ClassID
		quotas, err := cd.QuotaRepository.FindByScheduleID(ctx, tx, session.ScheduleID)
		if err != nil {

			return fmt.Errorf("fetch quotas failed: %w", err)
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
				return fmt.Errorf("quota not found for class %d", item.ClassID)
			}

			classData := dataQueue[item.ClassID]
			if len(classData) < item.Quantity {
				return fmt.Errorf("not enough ticket data for class %d", item.ClassID)
			}

			for i := 0; i < item.Quantity; i++ {
				data := classData[i]
				switch quota.Class.Type {
				case "passenger":
					if data.PassengerName == "" || data.IDType == "" || data.IDNumber == "" {

						return fmt.Errorf("missing passenger info for class %d", item.ClassID)
					}
					seat := fmt.Sprintf("%s%d", quota.Class.ClassAlias, quota.Capacity-quota.Quota+1)
					data.SeatNumber = &seat
				case "vehicle":
					if data.LicensePlate == nil || *data.LicensePlate == "" {

						return fmt.Errorf("missing license plate for vehicle class %d", item.ClassID)
					}
				default:

					return fmt.Errorf("unsupported ticket type or missing required fields for class %d", item.ClassID)
				}
				tickets = append(tickets, &domain.Ticket{
					TicketCode:      utils.GenerateTicketReferenceID(), // Unique ticket code
					BookingID:       &booking.ID,
					ClassID:         item.ClassID,
					Price:           quota.Price,
					Type:            quota.Class.Type,
					PassengerName:   data.PassengerName,
					PassengerAge:    data.PassengerAge,
					Address:         data.Address,
					PassengerGender: &data.PassengerGender,
					IDType:          &data.IDType,
					IDNumber:        &data.IDNumber,
					SeatNumber:      data.SeatNumber,
					LicensePlate:    data.LicensePlate,
					ScheduleID:      session.ScheduleID,
				})
				amounts += quota.Price
			}

			// Trim used data
			dataQueue[item.ClassID] = classData[item.Quantity:]
		}

		// Insert tickets
		if err := cd.TicketRepository.InsertBulk(ctx, tx, tickets); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create tickets: %w", err)
		}

		orderItems := make([]domain.OrderItem, len(tickets))
		for i, ticket := range tickets {
			orderItems[i] = client.TicketToItem(ticket)
		}

		payload := &domain.TransactionRequest{
			Method:        request.PaymentMethod,
			Amount:        int(amounts), // Convert to integer cents
			CustomerName:  booking.CustomerName,
			CustomerEmail: booking.Email,
			CustomerPhone: booking.PhoneNumber,
			MerchantRef:   booking.OrderID,
			OrderItems:    orderItems,
			CallbackUrl:   "https://example.com/callback",
			ReturnUrl:     "https://example.com/callback",
			ExpiredTime:   int(time.Now().Add(30 * time.Minute).Unix()),
		}

		payment, err := cd.TripayClient.CreatePayment(payload)
		if err != nil {
			return fmt.Errorf("create Tripay payment failed: %w", err)
		}

		booking.ReferenceNumber = &payment.Reference
		if err := cd.BookingRepository.Update(ctx, tx, booking); err != nil {
			return fmt.Errorf("failed to update booking with reference number: %w", err)
		}

		// Decrement quota usage
		for _, item := range session.ClaimItems {
			quota, ok := quotaByClass[item.ClassID]
			if !ok {
				return fmt.Errorf("quota not found for class %d", item.ClassID)
			}

			if quota.Quota < item.Quantity {
				return fmt.Errorf("not enough quota for class %d", item.ClassID)
			}

			quota.Quota -= item.Quantity

			if err := cd.QuotaRepository.Update(ctx, tx, quota); err != nil {
				return fmt.Errorf("failed to update quota usage: %w", err)
			}
		}

		subject := "Your Booking is Confirmed"
		htmlBody := templates.BookingInvoiceEmail(booking, payment)
		// cd.Mailer.SendAsync(booking.Email, subject, htmlBody)
		if err := cd.Mailer.Send(booking.Email, subject, htmlBody); err != nil {
			return fmt.Errorf("failed to send booking confirmation email: %w", err)
		}
		session.Status = enum.ClaimSessionSuccess.String()
		if err := cd.ClaimSessionRepository.Update(ctx, tx, session); err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("execute transaction: %w", err)
	}

	return &model.TESTReadClaimSessionDataEntryResponse{
		OrderID:   booking.OrderID,
		ExpiresAt: booking.ExpiresAt,
	}, nil
}

func (uc *ClaimSessionUsecase) CreateClaimSession(ctx context.Context, request *model.TESTWriteClaimSessionRequest) error {

	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimSession := &domain.ClaimSession{
			SessionID:  uuid.NewString(),
			ScheduleID: request.ScheduleID,
			Status:     enum.ClaimSessionPending.String(),
			ExpiresAt:  time.Now().Add(16 * time.Minute),
			ClaimItems: make([]domain.ClaimItem, len(request.Items)),
		}
		for i, item := range request.Items {
			claimSession.ClaimItems[i] = domain.ClaimItem{
				ClassID:  item.ClassID,
				Quantity: item.Quantity,
				Subtotal: item.Subtotal, // Assuming subtotal is provided in the request
			}
		}
		if err := uc.ClaimSessionRepository.Insert(ctx, tx, claimSession); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create claim session: %w", err)
		}

		return nil
	})
}

func (uc *ClaimSessionUsecase) UpdateClaimSession(ctx context.Context, request *model.UpdateClaimSessionRequest) error {

	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimSession, err := uc.ClaimSessionRepository.FindByID(ctx, tx, request.ID)
		if err != nil {
			return fmt.Errorf("failed to find quota: %w", err)
		}
		if claimSession == nil {
			return errs.ErrNotFound
		}

		claimSession.ScheduleID = request.ScheduleID
		claimSession.Status = request.Status
		claimSession.ExpiresAt = request.ExpiresAt
		claimSession.ClaimItems = make([]domain.ClaimItem, len(request.ClaimItems))
		for i, item := range request.ClaimItems {
			claimSession.ClaimItems[i] = domain.ClaimItem{
				ClassID:  item.ClassID,
				Quantity: item.Quantity,
				Subtotal: item.Subtotal, // Assuming subtotal is provided in the request
			}
		}

		if err := uc.ClaimSessionRepository.Update(ctx, tx, claimSession); err != nil {
			return fmt.Errorf("failed to create claim session: %w", err)
		}

		return nil
	})
}

func (uc *ClaimSessionUsecase) ListClaimSessions(ctx context.Context, limit, offset int, sort, search string) ([]*model.TESTReadClaimSessionResponse, int, error) {
	var err error
	var total int64
	var claimSessions []*domain.ClaimSession
	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.ClaimSessionRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count claim sessions: %w", err)
		}

		claimSessions, err = uc.ClaimSessionRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all claim sessions: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("execute transaction: %w", err)
	}
	responses := make([]*model.TESTReadClaimSessionResponse, len(claimSessions))
	for i, claimSession := range claimSessions {
		responses[i] = mapper.TESTClaimSessionToResponse(claimSession)
	}
	return responses, int(total), nil
}

func (uc *ClaimSessionUsecase) GetClaimSessionByID(ctx context.Context, id uint) (*model.TESTReadClaimSessionResponse, error) {
	var err error
	var claimSession *domain.ClaimSession
	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimSession, err = uc.ClaimSessionRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if claimSession == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get claim session by id: %w", err)
	}

	return mapper.TESTClaimSessionToResponse(claimSession), nil
}

func (uc *ClaimSessionUsecase) GetBySessionID(ctx context.Context, sessionUUID string) (*model.TESTReadClaimSessionResponse, error) {
	var err error
	var claimSession *domain.ClaimSession
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimSession, err = uc.ClaimSessionRepository.FindBySessionID(ctx, tx, sessionUUID)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if claimSession == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get claim session by session ID: %w", err)
	}
	return mapper.TESTClaimSessionToResponse(claimSession), nil
}

func (uc *ClaimSessionUsecase) DeleteClaimSession(ctx context.Context, id uint) error {

	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {

		claimSession, err := uc.ClaimSessionRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if claimSession == nil {
			return errs.ErrNotFound
		}

		if err := uc.ClaimSessionRepository.Delete(ctx, tx, claimSession); err != nil {
			return fmt.Errorf("failed to delete fare: %w", err)
		}
		return nil
	})
}

func (uc *ClaimSessionUsecase) DeleteExpiredClaimSession(ctx context.Context) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		expiredSessions, err := uc.ClaimSessionRepository.FindExpired(ctx, tx, 50)
		if err != nil {
			return fmt.Errorf("failed to find expired sessions: %w", err)
		}
		if len(expiredSessions) == 0 {
			return nil
		}

		if err := uc.ClaimSessionRepository.DeleteBulk(ctx, tx, expiredSessions); err != nil {
			return fmt.Errorf("failed to delete expired sessions: %w", err)
		}
		return nil
	})
}
