package usecase

import (
	"context"
	"testing"

	"eticket-api/internal/domain"
	"eticket-api/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func ticketUsecase(t *testing.T) (*TicketUsecase, *mocks.MockTicketRepository, *mocks.MockScheduleRepository, *mocks.MockQuotaRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	ticketRepo := mocks.NewMockTicketRepository(ctrl)
	scheduleRepo := mocks.NewMockScheduleRepository(ctrl)
	quotaRepo := mocks.NewMockQuotaRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := NewTicketUsecase(transactor, ticketRepo, scheduleRepo, quotaRepo)
	return uc, ticketRepo, scheduleRepo, quotaRepo, transactor
}

func TestTicketUsecase_CreateTicket(t *testing.T) {
	t.Parallel()
	uc, _, _, _, transactor := ticketUsecase(t)
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
			err := uc.CreateTicket(context.Background(), &domain.Ticket{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketUsecase_CheckIn(t *testing.T) {
	t.Parallel()
	uc, _, _, _, transactor := ticketUsecase(t)
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
			err := uc.CheckIn(context.Background(), 1)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
