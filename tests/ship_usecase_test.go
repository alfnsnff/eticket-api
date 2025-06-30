package tests

import (
	"context"
	"eticket-api/internal/domain"
	"eticket-api/internal/mocks"
	"eticket-api/internal/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// Helper for ShipUsecase
func shipUsecase(t *testing.T) (*usecase.ShipUsecase, *mocks.MockShipRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockShipRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewShipUsecase(transactor, repo)
	return uc, repo, transactor
}

func TestShipUsecase_CreateShip(t *testing.T) {
	t.Parallel()
	uc, _, transactor := shipUsecase(t)
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
			err := uc.CreateShip(context.Background(), &domain.Ship{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
