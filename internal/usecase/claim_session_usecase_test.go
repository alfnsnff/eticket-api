package usecase

import (
	"context"
	"testing"

	"eticket-api/internal/mocks"
	"eticket-api/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func claimSessionUsecase(t *testing.T) (*ClaimSessionUsecase, *mocks.MockClaimSessionRepository, *mocks.MockClaimItemRepository, *mocks.MockTicketRepository, *mocks.MockScheduleRepository, *mocks.MockBookingRepository, *mocks.MockQuotaRepository, *mocks.MockTripayClient, *mocks.MockMailer, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	claimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	claimItemRepo := mocks.NewMockClaimItemRepository(ctrl)
	ticketRepo := mocks.NewMockTicketRepository(ctrl)
	scheduleRepo := mocks.NewMockScheduleRepository(ctrl)
	bookingRepo := mocks.NewMockBookingRepository(ctrl)
	quotaRepo := mocks.NewMockQuotaRepository(ctrl)
	tripayClient := mocks.NewMockTripayClient(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := NewClaimSessionUsecase(transactor, claimSessionRepo, claimItemRepo, ticketRepo, scheduleRepo, bookingRepo, quotaRepo, tripayClient, mailer)
	return uc, claimSessionRepo, claimItemRepo, ticketRepo, scheduleRepo, bookingRepo, quotaRepo, tripayClient, mailer, transactor
}

func TestClaimSessionUsecase_CreateClaimSession(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, _, _, _, _, transactor := claimSessionUsecase(t)
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
			err := uc.CreateClaimSession(context.Background(), &model.TESTWriteClaimSessionRequest{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClaimSessionUsecase_DeleteExpiredClaimSession(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, _, _, _, _, transactor := claimSessionUsecase(t)
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
			err := uc.DeleteExpiredClaimSession(context.Background())
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Lakukan hal serupa untuk usecase lain (Class, Harbor, Quota, Schedule, Ticket, Booking, ClaimSession, ClaimItem, Payment)
// Copy helper dan test function di atas, ganti dependency dan method sesuai usecase/constructor masing-masing.
