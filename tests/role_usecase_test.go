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

// Helper for RoleUsecase
func roleUsecase(t *testing.T) (*usecase.RoleUsecase, *mocks.MockRoleRepository, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRoleRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := usecase.NewRoleUsecase(transactor, repo)
	return uc, repo, transactor
}

func TestRoleUsecase_CreateRole(t *testing.T) {
	t.Parallel()
	uc, _, transactor := roleUsecase(t)
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
			err := uc.CreateRole(context.Background(), &domain.Role{})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
