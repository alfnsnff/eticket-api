package tests

import (
	"context"
	"testing"

	"eticket-api/internal/domain"
	"eticket-api/internal/mocks"
	"eticket-api/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func scheduleUsecase(t *testing.T) (*usecase.ScheduleUsecase, *mocks.MockClaimSessionRepository, *mocks.MockClassRepository, *mocks.MockShipRepository, *mocks.MockScheduleRepository, *mocks.MockTicketRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	claimSessionRepo := mocks.NewMockClaimSessionRepository(ctrl)
	classRepo := mocks.NewMockClassRepository(ctrl)
	shipRepo := mocks.NewMockShipRepository(ctrl)
	scheduleRepo := mocks.NewMockScheduleRepository(ctrl)
	ticketRepo := mocks.NewMockTicketRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewScheduleUsecase(transactor, claimSessionRepo, classRepo, shipRepo, scheduleRepo, ticketRepo)
	return uc, claimSessionRepo, classRepo, shipRepo, scheduleRepo, ticketRepo, transactor
}

func TestScheduleUsecase_CreateSchedule(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, _, transactor := scheduleUsecase(t)
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
			err := uc.CreateSchedule(context.Background(), &domain.Schedule{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestScheduleUsecase_GetScheduleByID(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, _, transactor := scheduleUsecase(t)
	tests := []struct {
		name string
		mock func()
		res  *domain.Schedule
		err  error
	}{
		{
			name: "success",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
			},
			res: &domain.Schedule{},
			err: nil,
		},
		{
			name: "not found",
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
			_, err := uc.GetScheduleByID(context.Background(), 1)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
