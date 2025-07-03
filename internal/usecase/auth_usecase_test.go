package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"eticket-api/internal/common/token"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/mocks"
	"eticket-api/pkg/gotann"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var errInternalServErr = errors.New("internal server error")

// Helper untuk AuthUsecase
func authUsecase(t *testing.T) (*AuthUsecase, *mocks.MockRefreshTokenRepository, *mocks.MockUserRepository, *mocks.MockMailer, *mocks.MockTokenUtil, *mocks.MockTransactor) {
	t.Helper()
	ctrl := gomock.NewController(t)
	refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	tokenUtil := mocks.NewMockTokenUtil(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	uc := NewAuthUsecase(transactor, refreshRepo, userRepo, mailer, tokenUtil)
	return uc, refreshRepo, userRepo, mailer, tokenUtil, transactor
}

func TestAuthUsecase_Login(t *testing.T) {
	t.Parallel()
	uc, refreshRepo, userRepo, _, tokenUtil, transactor := authUsecase(t)

	userID := uint(1)
	hashedPassword, _ := utils.HashPassword("password123") // password = "password123"
	mockUser := &domain.User{
		ID:       userID,
		Username: "testuser",
		Password: hashedPassword, // Gunakan hashed password yang valid
	}

	tests := []struct {
		name     string
		request  *domain.Login
		mock     func()
		expected *domain.User
		hasError bool
	}{
		{
			name: "success",
			request: &domain.Login{
				Username: "testuser",
				Password: "password123",
			},
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil) // Execute the actual function
					},
				)
				userRepo.EXPECT().FindByUsername(gomock.Any(), gomock.Any(), "testuser").Return(mockUser, nil)
				tokenUtil.EXPECT().GenerateAccessToken(mockUser).Return("access_token", nil)
				tokenUtil.EXPECT().GenerateRefreshToken(mockUser).Return("refresh_token", nil)
				tokenUtil.EXPECT().ValidateToken("refresh_token").Return(&token.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ID:        uuid.New().String(),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					},
				}, nil)
				refreshRepo.EXPECT().InsertRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: mockUser,
			hasError: false,
		},
		{
			name: "user not found",
			request: &domain.Login{
				Username: "nonexistent",
				Password: "password123",
			},
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				userRepo.EXPECT().FindByUsername(gomock.Any(), gomock.Any(), "nonexistent").Return(nil, nil)
			},
			expected: nil,
			hasError: true,
		},
		{
			name: "wrong password",
			request: &domain.Login{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				userRepo.EXPECT().FindByUsername(gomock.Any(), gomock.Any(), "testuser").Return(mockUser, nil)
			},
			expected: nil,
			hasError: true,
		},
		{
			name: "database error",
			request: &domain.Login{
				Username: "testuser",
				Password: "password123",
			},
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				userRepo.EXPECT().FindByUsername(gomock.Any(), gomock.Any(), "testuser").Return(nil, errInternalServErr)
			},
			expected: nil,
			hasError: true,
		},
		{
			name: "token generation error",
			request: &domain.Login{
				Username: "testuser",
				Password: "password123",
			},
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				userRepo.EXPECT().FindByUsername(gomock.Any(), gomock.Any(), "testuser").Return(mockUser, nil)
				tokenUtil.EXPECT().GenerateAccessToken(mockUser).Return("", errors.New("token generation failed"))
			},
			expected: nil,
			hasError: true,
		},
		{
			name: "refresh token storage error",
			request: &domain.Login{
				Username: "testuser",
				Password: "password123",
			},
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				userRepo.EXPECT().FindByUsername(gomock.Any(), gomock.Any(), "testuser").Return(mockUser, nil)
				tokenUtil.EXPECT().GenerateAccessToken(mockUser).Return("access_token", nil)
				tokenUtil.EXPECT().GenerateRefreshToken(mockUser).Return("refresh_token", nil)
				tokenUtil.EXPECT().ValidateToken("refresh_token").Return(&token.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ID:        uuid.New().String(),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					},
				}, nil)
				refreshRepo.EXPECT().InsertRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(errInternalServErr)
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			user, accessToken, refreshToken, err := uc.Login(context.Background(), tc.request)

			if tc.hasError {
				require.Error(t, err)
				require.Nil(t, user)
				require.Empty(t, accessToken)
				require.Empty(t, refreshToken)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, user)
				require.NotEmpty(t, accessToken)
				require.NotEmpty(t, refreshToken)
			}
		})
	}
}

func TestAuthUsecase_RefreshToken(t *testing.T) {
	t.Parallel()
	uc, refreshRepo, userRepo, _, tokenUtil, transactor := authUsecase(t)

	userID := uint(1)
	refreshTokenID := uuid.New()
	mockUser := &domain.User{
		ID:       userID,
		Username: "testuser",
	}
	mockRefreshToken := &domain.RefreshToken{
		ID:        refreshTokenID,
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	tests := []struct {
		name     string
		token    string
		mock     func()
		hasError bool
	}{
		{
			name:  "success",
			token: "valid_refresh_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("valid_refresh_token").Return(&token.Claims{
					User: mockUser,
					RegisteredClaims: jwt.RegisteredClaims{
						ID: refreshTokenID.String(),
					},
				}, nil)
				refreshRepo.EXPECT().FindRefreshTokenByIDAndStatus(gomock.Any(), gomock.Any(), refreshTokenID.String(), false).Return(mockRefreshToken, nil)
				userRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), userID).Return(mockUser, nil)
				tokenUtil.EXPECT().GenerateAccessToken(mockUser).Return("new_access_token", nil)
			},
			hasError: false,
		},
		{
			name:  "invalid token",
			token: "invalid_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("invalid_token").Return(nil, errors.New("invalid token"))
			},
			hasError: true,
		},
		{
			name:  "revoked token",
			token: "revoked_token",
			mock: func() {
				revokedToken := &domain.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					Revoked:   true,
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("revoked_token").Return(&token.Claims{
					User: mockUser,
					RegisteredClaims: jwt.RegisteredClaims{
						ID: refreshTokenID.String(),
					},
				}, nil)
				refreshRepo.EXPECT().FindRefreshTokenByIDAndStatus(gomock.Any(), gomock.Any(), refreshTokenID.String(), false).Return(revokedToken, nil)
			},
			hasError: true,
		},
		{
			name:  "expired token",
			token: "expired_token",
			mock: func() {
				expiredToken := &domain.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					Revoked:   false,
					ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired
				}
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("expired_token").Return(&token.Claims{
					User: mockUser,
					RegisteredClaims: jwt.RegisteredClaims{
						ID: refreshTokenID.String(),
					},
				}, nil)
				refreshRepo.EXPECT().FindRefreshTokenByIDAndStatus(gomock.Any(), gomock.Any(), refreshTokenID.String(), false).Return(expiredToken, nil)
			},
			hasError: true,
		},
		{
			name:  "user not found",
			token: "valid_token_no_user",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("valid_token_no_user").Return(&token.Claims{
					User: mockUser,
					RegisteredClaims: jwt.RegisteredClaims{
						ID: refreshTokenID.String(),
					},
				}, nil)
				refreshRepo.EXPECT().FindRefreshTokenByIDAndStatus(gomock.Any(), gomock.Any(), refreshTokenID.String(), false).Return(mockRefreshToken, nil)
				userRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), userID).Return(nil, nil)
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			newToken, err := uc.RefreshToken(context.Background(), tc.token)

			if tc.hasError {
				require.Error(t, err)
				require.Empty(t, newToken)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, newToken)
			}
		})
	}
}

func TestAuthUsecase_Logout(t *testing.T) {
	t.Parallel()
	uc, refreshRepo, _, _, tokenUtil, transactor := authUsecase(t)

	tokenID := uuid.New()

	tests := []struct {
		name     string
		token    string
		mock     func()
		hasError bool
	}{
		{
			name:  "success",
			token: "valid_refresh_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("valid_refresh_token").Return(&token.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ID: tokenID.String(),
					},
				}, nil)
				refreshRepo.EXPECT().RevokeRefreshTokenByID(gomock.Any(), gomock.Any(), tokenID).Return(nil)
			},
			hasError: false,
		},
		{
			name:  "invalid token",
			token: "invalid_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("invalid_token").Return(nil, errors.New("invalid token"))
			},
			hasError: true,
		},
		{
			name:  "invalid token ID format",
			token: "invalid_id_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("invalid_id_token").Return(&token.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ID: "invalid-uuid-format",
					},
				}, nil)
			},
			hasError: true,
		},
		{
			name:  "revoke error",
			token: "revoke_error_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("revoke_error_token").Return(&token.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ID: tokenID.String(),
					},
				}, nil)
				refreshRepo.EXPECT().RevokeRefreshTokenByID(gomock.Any(), gomock.Any(), tokenID).Return(errInternalServErr)
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := uc.Logout(context.Background(), tc.token)

			if tc.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_Me(t *testing.T) {
	t.Parallel()
	uc, _, userRepo, _, tokenUtil, transactor := authUsecase(t)

	userID := uint(1)
	mockUser := &domain.User{
		ID:       userID,
		Username: "testuser",
	}

	tests := []struct {
		name     string
		token    string
		mock     func()
		expected *domain.User
		hasError bool
	}{
		{
			name:  "success",
			token: "valid_access_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("valid_access_token").Return(&token.Claims{
					User: mockUser,
				}, nil)
				userRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), userID).Return(mockUser, nil)
			},
			expected: mockUser,
			hasError: false,
		},
		{
			name:  "invalid token",
			token: "invalid_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("invalid_token").Return(nil, errors.New("invalid token"))
			},
			expected: nil,
			hasError: true,
		},
		{
			name:  "user not found",
			token: "valid_token_no_user",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("valid_token_no_user").Return(&token.Claims{
					User: mockUser,
				}, nil)
				userRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), userID).Return(nil, nil)
			},
			expected: nil,
			hasError: true,
		},
		{
			name:  "database error",
			token: "db_error_token",
			mock: func() {
				transactor.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx gotann.Transaction) error) error {
						return fn(nil)
					},
				)
				tokenUtil.EXPECT().ValidateToken("db_error_token").Return(&token.Claims{
					User: mockUser,
				}, nil)
				userRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), userID).Return(nil, errInternalServErr)
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			user, err := uc.Me(context.Background(), tc.token)

			if tc.hasError {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, user)
			}
		})
	}
}
