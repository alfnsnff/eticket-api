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

func bookingUsecase(t *testing.T) (*usecase.BookingUsecase, *mocks.MockBookingRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockBookingRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewBookingUsecase(transactor, repo)
	return uc, repo, transactor
}

func TestBookingUsecase_CreateBooking(t *testing.T) {
	t.Parallel()
	uc, _, transactor := bookingUsecase(t)
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
			err := uc.CreateBooking(context.Background(), &domain.Booking{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
