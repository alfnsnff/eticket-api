package usecase

import (
	"context"
	"testing"

	"eticket-api/internal/domain"
	"eticket-api/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func harborUsecase(t *testing.T) (*HarborUsecase, *mocks.MockHarborRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockHarborRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := NewHarborUsecase(transactor, repo)
	return uc, repo, transactor
}

func TestHarborUsecase_CreateHarbor(t *testing.T) {
	t.Parallel()
	uc, _, transactor := harborUsecase(t)
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
			err := uc.CreateHarbor(context.Background(), &domain.Harbor{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHarborUsecase_GetHarborByID(t *testing.T) {
	t.Parallel()
	uc, _, transactor := harborUsecase(t)
	tests := []struct {
		name string
		mock func()
		res  *domain.Harbor
		err  error
	}{
		{
			name: "success",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
			},
			res: &domain.Harbor{},
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
			_, err := uc.GetHarborByID(context.Background(), 1)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
