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

func classUsecase(t *testing.T) (*usecase.ClassUsecase, *mocks.MockClassRepository, *mocks.MockTransactor) {
    t.Helper()
    ctrl := gomock.NewController(t)
    repo := mocks.NewMockClassRepository(ctrl)
    transactor := mocks.NewMockTransactor(ctrl)
    uc := usecase.NewClassUsecase(transactor, repo)
    return uc, repo, transactor
}

func TestClassUsecase_CreateClass(t *testing.T) {
    t.Parallel()
    uc, _, transactor := classUsecase(t)
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
            err := uc.CreateClass(context.Background(), &domain.Class{})
            if tc.err != nil {
                require.ErrorIs(t, err, tc.err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}

func TestClassUsecase_GetClassByID(t *testing.T) {
    t.Parallel()
    uc, _, transactor := classUsecase(t)
    tests := []struct {
        name string
        mock func()
        res  *domain.Class
        err  error
    }{
        {
            name: "success",
            mock: func() {
                transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
            },
            res: &domain.Class{},
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
            _, err := uc.GetClassByID(context.Background(), 1)
            if tc.err != nil {
                require.ErrorIs(t, err, tc.err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}