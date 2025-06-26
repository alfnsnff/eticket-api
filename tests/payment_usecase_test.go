package tests

import (
	"context"
	"testing"
	"time"

	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/domain/mocks"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"

	"github.com/glebarez/sqlite"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupPaymentTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory sqlite: " + err.Error())
	}
	db.AutoMigrate(&domain.ClaimSession{}, &domain.Booking{}, &domain.Ticket{})
	return db
}

func TestListPaymentChannels_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()
	expectedChannels := []model.ReadPaymentChannelResponse{
		{
			Group: "Virtual Account",
			Code:  "BRIVA",
			Name:  "BRI Virtual Account",
			Type:  "virtual_account",
			FeeCustomer: model.Fee{
				Flat:    4000,
				Percent: 0,
			},
			MinimumAmount: 10000,
			MaximumAmount: 1000000000,
			IconURL:       "https://tripay.co.id/images/payment/briva.png",
			Active:        true,
		},
	}

	mockTripayClient.EXPECT().
		GetPaymentChannels().
		Return(expectedChannels, nil).
		Times(1)

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	channels, err := usecase.ListPaymentChannels(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, channels)
	assert.Len(t, channels, 1)
	assert.Equal(t, "BRIVA", channels[0].Code)
	assert.Equal(t, "BRI Virtual Account", channels[0].Name)
}

func TestListPaymentChannels_ClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	mockTripayClient.EXPECT().
		GetPaymentChannels().
		Return(nil, assert.AnError).
		Times(1)

	channels, err := usecase.ListPaymentChannels(context.Background())
	assert.Error(t, err)
	assert.Nil(t, channels)
}

func TestGetTransactionDetail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	expectedDetail := &model.ReadTransactionResponse{
		Reference:   "test-reference",
		Status:      "PAID",
		Amount:      100000,
		CheckoutUrl: "https://tripay.co.id/checkout/test",
	}

	mockTripayClient.EXPECT().
		GetTransactionDetail("test-reference").
		Return(expectedDetail, nil).
		Times(1)

	detail, err := usecase.GetTransactionDetail(context.Background(), "test-reference")
	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "test-reference", detail.Reference)
	assert.Equal(t, "PAID", detail.Status)
}

func TestGetTransactionDetail_ClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	mockTripayClient.EXPECT().
		GetTransactionDetail("test-reference").
		Return(nil, assert.AnError).
		Times(1)

	detail, err := usecase.GetTransactionDetail(context.Background(), "test-reference")
	assert.Error(t, err)
	assert.Nil(t, detail)
}

func TestCreatePayment_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WritePaymentRequest{
		OrderID:       "test-order-123",
		PaymentMethod: "BRIVA",
	}

	mockClaimSessionRepo.EXPECT().
		FindBySessionID(gomock.Any(), "test-session-id").
		Return(nil, nil).Times(1)

	result, err := usecase.CreatePayment(context.Background(), req, "test-session-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestCreatePayment_SessionExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WritePaymentRequest{
		OrderID:       "test-order-123",
		PaymentMethod: "BRIVA",
	}

	expiredSession := &domain.ClaimSession{
		ID:        1,
		SessionID: "test-session-id",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		Status:    "PAYMENT_PENDING",
	}

	mockClaimSessionRepo.EXPECT().
		FindBySessionID(gomock.Any(), "test-session-id").
		Return(expiredSession, nil).Times(1)

	_, err := usecase.CreatePayment(context.Background(), req, "test-session-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.ErrExpired.Error())
}

func TestCreatePayment_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WritePaymentRequest{
		OrderID:       "test-order-123",
		PaymentMethod: "BRIVA",
	}
	validSession := &domain.ClaimSession{
		ID:        1,
		SessionID: "test-session-id",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Status:    "PAYMENT_PENDING",
	}

	mockClaimSessionRepo.EXPECT().
		FindBySessionID(gomock.Any(), "test-session-id").
		Return(validSession, nil).Times(1)

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(nil, nil).Times(1)

	_, err := usecase.CreatePayment(context.Background(), req, "test-session-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestCreatePayment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WritePaymentRequest{
		OrderID:       "test-order-123",
		PaymentMethod: "BRIVA",
	}

	validSession := &domain.ClaimSession{
		ID:        1,
		SessionID: "test-session-id",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Status:    "PAYMENT_PENDING",
	}

	booking := &domain.Booking{
		ID:           1,
		OrderID:      &req.OrderID,
		CustomerName: "Test Customer",
		Email:        "test@example.com",
	}

	claimSessionID := uint(1)
	tickets := []*domain.Ticket{
		{
			ID:             1,
			ClaimSessionID: &claimSessionID,
			PassengerName:  stringPtr("Test Passenger"),
			Price:          100000,
		},
	}

	expectedResponse := model.ReadTransactionResponse{
		Reference:   "TRIPAY-REF-123",
		Status:      "UNPAID",
		CheckoutUrl: "https://tripay.co.id/checkout/test",
		Amount:      100000,
	}

	mockClaimSessionRepo.EXPECT().
		FindBySessionID(gomock.Any(), "test-session-id").
		Return(validSession, nil).Times(1)

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(booking, nil).Times(1)

	mockTicketRepo.EXPECT().
		FindByBookingID(gomock.Any(), uint(1)).
		Return(tickets, nil).Times(1)
	mockTripayClient.EXPECT().
		CreatePayment(gomock.Any()).
		DoAndReturn(func(payload *model.WriteTransactionRequest) (model.ReadTransactionResponse, error) {
			assert.Equal(t, "BRIVA", payload.Method)
			assert.Equal(t, "test-order-123", payload.MerchantRef)
			assert.Equal(t, 100000, payload.Amount)
			assert.Equal(t, "Test Customer", payload.CustomerName)
			assert.Equal(t, "test@example.com", payload.CustomerEmail)
			return expectedResponse, nil
		}).Times(1)

	mockBookingRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, booking *domain.Booking) error {
			assert.Equal(t, "TRIPAY-REF-123", *booking.ReferenceNumber)
			return nil
		}).Times(1)

	mockClaimSessionRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, sess *domain.ClaimSession) error {
			assert.Equal(t, "TRANSACTION_PENDING", sess.Status)
			return nil
		}).Times(1)

	result, err := usecase.CreatePayment(context.Background(), req, "test-session-id")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TRIPAY-REF-123", result.Reference)
	assert.Equal(t, "UNPAID", result.Status)
}

func TestHandleCallback_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WriteCallbackRequest{
		Reference:     "test-reference",
		MerchantRef:   "test-order-123",
		Status:        "PAID",
		Amount:        100000,
		PaymentMethod: "BRIVA",
		Signature:     "test-signature",
	}

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(nil, nil).Times(1)

	err := usecase.HandleCallback(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestHandleCallback_NoTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WriteCallbackRequest{
		Reference:     "test-reference",
		MerchantRef:   "test-order-123",
		Status:        "PAID",
		Amount:        100000,
		PaymentMethod: "BRIVA",
		Signature:     "test-signature",
	}

	booking := &domain.Booking{
		ID:           1,
		OrderID:      &req.MerchantRef,
		CustomerName: "Test Customer",
		Email:        "test@example.com",
	}

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(booking, nil).Times(1)

	mockTicketRepo.EXPECT().
		FindByBookingID(gomock.Any(), uint(1)).
		Return([]*domain.Ticket{}, nil).Times(1)

	err := usecase.HandleCallback(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
}

func TestHandleCallback_SuccessfulPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WriteCallbackRequest{
		Reference:     "test-reference",
		MerchantRef:   "test-order-123",
		Status:        "PAID",
		Amount:        100000,
		PaymentMethod: "BRIVA",
		Signature:     "test-signature",
	}

	booking := &domain.Booking{
		ID:           1,
		OrderID:      &req.MerchantRef,
		CustomerName: "Test Customer",
		Email:        "test@example.com",
	}

	claimSessionID := uint(1)
	tickets := []*domain.Ticket{
		{
			ID:             1,
			ClaimSessionID: &claimSessionID,
			PassengerName:  stringPtr("Test Passenger"),
			Price:          100000,
		},
	}

	session := &domain.ClaimSession{
		ID:        1,
		SessionID: "test-session-id",
		Status:    "TRANSACTION_PENDING",
	}

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(booking, nil).Times(1)

	mockTicketRepo.EXPECT().
		FindByBookingID(gomock.Any(), uint(1)).
		Return(tickets, nil).Times(1)

	mockClaimSessionRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(session, nil).Times(1)
	mockClaimSessionRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, sess *domain.ClaimSession) error {
			assert.Equal(t, "SUCCESS", sess.Status)
			return nil
		}).Times(1)

	mockMailer.EXPECT().
		Send("test@example.com", gomock.Any(), gomock.Any()).DoAndReturn(func(toEmail, subject, body string) error {
		assert.Contains(t, subject, "Booking is Confirmed")
		return nil
	}).Times(1)

	err := usecase.HandleCallback(context.Background(), req)
	assert.NoError(t, err)
}

func TestHandleCallback_FailedPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WriteCallbackRequest{
		Reference: "test-reference", MerchantRef: "test-order-123", Status: "FAILED",
		Amount:        100000,
		PaymentMethod: "BRIVA",
		Signature:     "test-signature",
	}

	booking := &domain.Booking{
		ID:           1,
		OrderID:      &req.MerchantRef,
		CustomerName: "Test Customer",
		Email:        "test@example.com",
	}

	claimSessionID := uint(1)
	tickets := []*domain.Ticket{
		{
			ID:             1,
			ClaimSessionID: &claimSessionID,
			PassengerName:  stringPtr("Test Passenger"),
			Price:          100000,
		},
	}

	session := &domain.ClaimSession{
		ID:        1,
		SessionID: "test-session-id",
		Status:    "TRANSACTION_PENDING",
	}

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(booking, nil).Times(1)

	mockTicketRepo.EXPECT().
		FindByBookingID(gomock.Any(), uint(1)).
		Return(tickets, nil).Times(1)

	mockClaimSessionRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(session, nil).Times(1)
	mockClaimSessionRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tx *gorm.DB, sess *domain.ClaimSession) error {
			assert.Equal(t, "FAILED", sess.Status)
			return nil
		}).Times(1)

	mockMailer.EXPECT().
		Send("test@example.com", gomock.Any(), gomock.Any()).
		DoAndReturn(func(toEmail, subject, body string) error {
			assert.Contains(t, subject, "Payment Failed")
			return nil
		}).Times(1)

	err := usecase.HandleCallback(context.Background(), req)
	assert.NoError(t, err)
}

func TestHandleCallback_UnknownStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripayClient := mocks.NewMockTripayClient(ctrl)
	mockClaimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockTicketRepo := mocks.NewMockTicketRepository(ctrl)
	mockMailer := mocks.NewMockMailer(ctrl)
	db := setupPaymentTestDatabase()

	usecase := usecase.NewPaymentUsecase(db, mockTripayClient, mockClaimSessionRepo, mockBookingRepo, mockTicketRepo, mockMailer)

	req := &model.WriteCallbackRequest{
		Reference: "test-reference", MerchantRef: "test-order-123", Status: "UNKNOWN_STATUS",
		Amount:        100000,
		PaymentMethod: "BRIVA",
		Signature:     "test-signature",
	}

	booking := &domain.Booking{
		ID:           1,
		OrderID:      &req.MerchantRef,
		CustomerName: "Test Customer",
		Email:        "test@example.com",
	}

	claimSessionID := uint(1)
	tickets := []*domain.Ticket{
		{
			ID:             1,
			ClaimSessionID: &claimSessionID,
			PassengerName:  stringPtr("Test Passenger"),
			Price:          100000,
		},
	}

	session := &domain.ClaimSession{
		ID:        1,
		SessionID: "test-session-id",
		Status:    "TRANSACTION_PENDING",
	}

	mockBookingRepo.EXPECT().
		FindByOrderID(gomock.Any(), "test-order-123").
		Return(booking, nil).Times(1)

	mockTicketRepo.EXPECT().
		FindByBookingID(gomock.Any(), uint(1)).
		Return(tickets, nil).Times(1)

	mockClaimSessionRepo.EXPECT().
		FindByID(gomock.Any(), uint(1)).
		Return(session, nil).Times(1)

	err := usecase.HandleCallback(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown payment status: UNKNOWN_STATUS")
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
