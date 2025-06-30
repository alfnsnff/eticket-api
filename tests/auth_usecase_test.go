package tests

import (
	"context"
	"errors"
	"testing"

	"eticket-api/internal/domain"
	"eticket-api/internal/mocks" // <--- pastikan ini path mock hasil generate GoMock
	"eticket-api/internal/usecase"
	"eticket-api/pkg/gotann"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errInternalServErr = errors.New("internal server error")

// Helper untuk AuthUsecase
func authUsecase(t *testing.T) (*usecase.AuthUsecase, *mocks.MockRefreshTokenRepository, *mocks.MockUserRepository, *mocks.MockMailer, *mocks.MockTokenUtil, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	tokenUtil := mocks.NewMockTokenUtil(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewAuthUsecase(transactor, refreshRepo, userRepo, mailer, tokenUtil)
	return uc, refreshRepo, userRepo, mailer, tokenUtil, transactor
}

func TestAuthUsecase_Login(t *testing.T) {
	t.Parallel()
	uc, _, _, _, _, transactor := authUsecase(t)

	tests := []struct {
		name string
		mock func()
		res  *domain.User
		err  error
	}{
		{
			name: "invalid credentials",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return errors.New("invalid credentials")
					},
				)
			},
			res: nil,
			err: errors.New("login failed: invalid credentials"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			res, _, _, err := uc.Login(context.Background(), &domain.Login{})
			require.Equal(t, tc.res, res)
			require.Error(t, err)
		})
	}
}
