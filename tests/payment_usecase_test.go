package tests

import (
	"context"
	"testing"

	"eticket-api/internal/domain"
	"eticket-api/internal/mocks"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func paymentUsecase(t *testing.T) (*usecase.PaymentUsecase, *mocks.MockTripayClient, *mocks.MockBookingRepository, *mocks.MockTicketRepository, *mocks.MockQuotaRepository, *mocks.MockMailer, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	tripayClient := mocks.NewMockTripayClient(ctrl)
	bookingRepo := mocks.NewMockBookingRepository(ctrl)
	ticketRepo := mocks.NewMockTicketRepository(ctrl)
	quotaRepo := mocks.NewMockQuotaRepository(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewPaymentUsecase(transactor, tripayClient, bookingRepo, ticketRepo, quotaRepo, mailer)
	return uc, tripayClient, bookingRepo, ticketRepo, quotaRepo, mailer, transactor
}

func TestPaymentUsecase_CreatePayment(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, _, transactor := paymentUsecase(t)
	tests := []struct {
		name string
		mock func()
		res  *domain.Transaction
		err  error
	}{
		{
			name: "success",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
			},
			res: &domain.Transaction{},
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errInternalServErr)
			},
			res: nil,
			err: errInternalServErr,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			_, err := uc.CreatePayment(context.Background(), &model.WritePaymentRequest{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPaymentUsecase_HandleCallback(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, _, transactor := paymentUsecase(t)
	tests := []struct {
		name string
		mock func()
		err  error
	}{
		{
			name: "success",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errInternalServErr)
			},
			err: errInternalServErr,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := uc.HandleCallback(context.Background(), &domain.Callback{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
