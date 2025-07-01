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

func userUsecase(t *testing.T) (*usecase.UserUsecase, *mocks.MockUserRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockUserRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewUserUsecase(transactor, repo)
	return uc, repo, transactor
}

func TestUserUsecase_CreateUser(t *testing.T) {
	t.Parallel()
	uc, _, transactor := userUsecase(t)
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
			err := uc.CreateUser(context.Background(), &domain.User{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
