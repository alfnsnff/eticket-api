package tests

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"

// 	errs "eticket-api/internal/common/errors"
// 	"eticket-api/internal/common/token"
// 	"eticket-api/internal/common/utils"
// 	"eticket-api/internal/domain"
// 	"eticket-api/internal/domain/mocks"
// 	"eticket-api/internal/model"
// 	"eticket-api/internal/usecase"

// 	"github.com/glebarez/sqlite"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/gorm"
// )

// func setupAuthTestDatabase() *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_fk=1"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect to in-memory sqlite: " + err.Error())
// 	}
// 	db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.RefreshToken{}, &domain.PasswordReset{})
// 	return db
// }

// func TestAuthUsecase_Login_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	req := &model.WriteLoginRequest{
// 		Username: "testuser",
// 		Password: "testpassword",
// 	}

// 	passwordHash, _ := utils.HashPassword("testpassword")

// 	user := &domain.User{
// 		ID:       1,
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 		Password: passwordHash, // Mock bcrypt hash
// 		Role: domain.Role{
// 			ID:       1,
// 			RoleName: "admin",
// 		},
// 	}
// 	// Mock token generation
// 	mockTokenUtil.EXPECT().
// 		GenerateAccessToken(user).
// 		Return("access-token", nil).Times(1)
// 	mockTokenUtil.EXPECT().
// 		GenerateRefreshToken(user).
// 		Return("refresh-token", nil).Times(1)

// 	// Mock token validation (happens after generation to extract claims)
// 	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
// 	expirationTime := testTime.Add(time.Hour)

// 	mockTokenUtil.EXPECT().
// 		ValidateToken("refresh-token").
// 		Return(&token.Claims{
// 			User: user,
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(expirationTime),
// 				IssuedAt:  jwt.NewNumericDate(testTime),
// 				Issuer:    "eticket-api",
// 				Subject:   fmt.Sprintf("%d", user.ID),
// 				ID:        "550e8400-e29b-41d4-a716-446655440000",
// 			},
// 		}, nil).Times(1)

// 	mockUserRepo.EXPECT().
// 		FindByUsername(gomock.Any(), req.Username).
// 		Return(user, nil).Times(1)

// 	mockAuthRepo.EXPECT().
// 		InsertRefreshToken(gomock.Any(), gomock.Any()).
// 		DoAndReturn(func(db *gorm.DB, rt *domain.RefreshToken) error {
// 			assert.Equal(t, user.ID, rt.UserID)
// 			assert.False(t, rt.Revoked)
// 			return nil
// 		}).Times(1)

// 	accessToken, refreshToken, err := usecase.Login(context.Background(), req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "access-token", accessToken)
// 	assert.Equal(t, "refresh-token", refreshToken)
// }

// func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	req := &model.WriteLoginRequest{
// 		Username: "nonexistent",
// 		Password: "testpassword",
// 	}

// 	mockUserRepo.EXPECT().
// 		FindByUsername(gomock.Any(), req.Username).
// 		Return(nil, nil).Times(1)

// 	accessToken, refreshToken, err := usecase.Login(context.Background(), req)

// 	assert.Error(t, err)
// 	assert.Empty(t, accessToken)
// 	assert.Empty(t, refreshToken)
// 	assert.Contains(t, err.Error(), "invalid credentials")
// }

// func TestAuthUsecase_Login_RepositoryError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	req := &model.WriteLoginRequest{
// 		Username: "testuser",
// 		Password: "testpassword",
// 	}

// 	mockUserRepo.EXPECT().
// 		FindByUsername(gomock.Any(), req.Username).
// 		Return(nil, errors.New("database error")).Times(1)

// 	accessToken, refreshToken, err := usecase.Login(context.Background(), req)

// 	assert.Error(t, err)
// 	assert.Empty(t, accessToken)
// 	assert.Empty(t, refreshToken)
// 	assert.Contains(t, err.Error(), "failed to retrieve user")
// }

// func TestAuthUsecase_RefreshToken_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)
// 	refreshToken := "valid-refresh-token"
// 	userID := uint(1)
// 	tokenID := "550e8400-e29b-41d4-a716-446655440000"

// 	user := &domain.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 		Role: domain.Role{
// 			ID:       1,
// 			RoleName: "admin",
// 		},
// 	}

// 	refreshTokenDomain := &domain.RefreshToken{
// 		ID:        uuid.MustParse(tokenID),
// 		UserID:    userID,
// 		Revoked:   false,
// 		ExpiresAt: time.Now().Add(time.Hour),
// 	}

// 	// Better approach - use consistent test times
// 	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
// 	expirationTime := testTime.Add(time.Hour)

// 	mockTokenUtil.EXPECT().
// 		ValidateToken(refreshToken).
// 		Return(&token.Claims{
// 			User: user,
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(expirationTime),
// 				IssuedAt:  jwt.NewNumericDate(testTime),
// 				Issuer:    "eticket-api",
// 				Subject:   fmt.Sprintf("%d", user.ID),
// 				ID:        tokenID,
// 			},
// 		}, nil).Times(1)
// 	mockTokenUtil.EXPECT().
// 		GenerateAccessToken(user).
// 		Return("new-access-token", nil).Times(1)

// 	mockAuthRepo.EXPECT().
// 		FindRefreshTokenByIDAndStatus(gomock.Any(), tokenID, false).
// 		Return(refreshTokenDomain, nil).Times(1)

// 	mockUserRepo.EXPECT().
// 		FindByID(gomock.Any(), userID).
// 		Return(user, nil).Times(1)

// 	newAccessToken, err := usecase.RefreshToken(context.Background(), refreshToken)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "new-access-token", newAccessToken)
// }

// func TestAuthUsecase_RefreshToken_InvalidToken(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	invalidToken := "invalid-token"

// 	mockTokenUtil.EXPECT().
// 		ValidateToken(invalidToken).
// 		Return(nil, errors.New("invalid token")).Times(1)

// 	newAccessToken, err := usecase.RefreshToken(context.Background(), invalidToken)

// 	assert.Error(t, err)
// 	assert.Empty(t, newAccessToken)
// 	assert.Contains(t, err.Error(), "invalid refresh token")
// }
// func TestAuthUsecase_Logout_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	refreshToken := "valid-refresh-token"

// 	// Define the token ID that will be used consistently
// 	tokenIDStr := "550e8400-e29b-41d4-a716-446655440000"
// 	tokenID, _ := uuid.Parse(tokenIDStr)

// 	// Create a test user
// 	testUser := &domain.User{
// 		ID:    1,
// 		Email: "test@example.com",
// 		// ... other required fields
// 	}

// 	// Better approach - use consistent test times
// 	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
// 	expirationTime := testTime.Add(time.Hour)

// 	mockTokenUtil.EXPECT().
// 		ValidateToken(refreshToken).
// 		Return(&token.Claims{
// 			User: testUser, // ✅ Now properly defined
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(expirationTime),
// 				IssuedAt:  jwt.NewNumericDate(testTime),
// 				Issuer:    "eticket-api",
// 				Subject:   fmt.Sprintf("%d", testUser.ID),
// 				ID:        tokenIDStr, // ✅ Use consistent ID
// 			},
// 		}, nil).Times(1)

// 	// The repository method receives a transaction, not the original DB
// 	mockAuthRepo.EXPECT().
// 		RevokeRefreshTokenByID(gomock.Any(), tokenID). // This now matches the parsed token ID
// 		Return(nil).
// 		Times(1)

// 	err := usecase.Logout(context.Background(), refreshToken)

// 	assert.NoError(t, err)
// }

// func TestAuthUsecase_Logout_InvalidToken(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	invalidToken := "invalid-token"

// 	mockTokenUtil.EXPECT().
// 		ValidateToken(invalidToken).
// 		Return(nil, errors.New("invalid token")).Times(1)

// 	err := usecase.Logout(context.Background(), invalidToken)

// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "invalid refresh token")
// }

// func TestAuthUsecase_Me_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	accessToken := "valid-access-token"
// 	userID := uint(1)

// 	user := &domain.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 		FullName: "Test User",
// 		Role: domain.Role{
// 			ID:       1,
// 			RoleName: "admin",
// 		},
// 	}

// 	// Better approach - use consistent test times
// 	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
// 	expirationTime := testTime.Add(time.Hour)

// 	mockTokenUtil.EXPECT().
// 		ValidateToken(accessToken).
// 		Return(&token.Claims{
// 			User: user,
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(expirationTime),
// 				IssuedAt:  jwt.NewNumericDate(testTime),
// 				Issuer:    "eticket-api",
// 				Subject:   fmt.Sprintf("%d", user.ID),
// 				ID:        "550e8400-e29b-41d4-a716-446655440000",
// 			},
// 		}, nil).Times(1)

// 	mockUserRepo.EXPECT().
// 		FindByID(gomock.Any(), userID).
// 		Return(user, nil).Times(1)

// 	response, err := usecase.Me(context.Background(), accessToken)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, response)
// 	assert.Equal(t, user.ID, response.ID)
// 	assert.Equal(t, user.Username, response.Username)
// 	assert.Equal(t, user.Email, response.Email)
// 	assert.Equal(t, user.FullName, response.FullName)
// }

// func TestAuthUsecase_RequestPasswordReset_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	email := "test@example.com"
// 	user := &domain.User{
// 		ID:       1,
// 		Username: "testuser",
// 		Email:    email,
// 		FullName: "Test User",
// 	}

// 	mockUserRepo.EXPECT().
// 		FindByEmail(gomock.Any(), email).
// 		Return(user, nil).Times(1)
// 	mockAuthRepo.EXPECT().
// 		InsertPasswordReset(gomock.Any(), gomock.Any()).
// 		DoAndReturn(func(db *gorm.DB, pr *domain.PasswordReset) error {
// 			assert.Equal(t, user.ID, pr.UserID)
// 			assert.NotEmpty(t, pr.Token)
// 			assert.True(t, pr.ExpiresAt.After(time.Now()))
// 			return nil
// 		}).Times(1)

// 	mockMailer.EXPECT().
// 		Send(email, "Password Reset", gomock.Any()).
// 		Return(nil).Times(1)

// 	err := usecase.RequestPasswordReset(context.Background(), email)

// 	assert.NoError(t, err)
// }

// func TestAuthUsecase_RequestPasswordReset_UserNotFound(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	email := "nonexistent@example.com"

// 	mockUserRepo.EXPECT().
// 		FindByEmail(gomock.Any(), email).
// 		Return(nil, nil).Times(1)

// 	err := usecase.RequestPasswordReset(context.Background(), email)

// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), errs.ErrNotFound.Error())
// }

// func TestAuthUsecase_ResetPassword_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockMailer := mocks.NewMockMailer(ctrl)
// 	mockTokenUtil := mocks.NewMockTokenUtil(ctrl)
// 	db := setupAuthTestDatabase()

// 	usecase := usecase.NewAuthUsecase(db, mockAuthRepo, mockUserRepo, mockMailer, mockTokenUtil)

// 	token := "valid-reset-token"
// 	newPassword := "newpassword123"
// 	userID := uint(1)

// 	user := &domain.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 		Password: "oldhashedpassword",
// 	}

// 	passwordReset := &domain.PasswordReset{
// 		UserID:    userID,
// 		Token:     token,
// 		Issued:    false,
// 		ExpiresAt: time.Now().Add(time.Hour),
// 	}

// 	mockAuthRepo.EXPECT().
// 		FindPasswordResetByTokenAndStatus(gomock.Any(), token, false).
// 		Return(passwordReset, nil).Times(1)

// 	mockUserRepo.EXPECT().
// 		FindByID(gomock.Any(), userID).
// 		Return(user, nil).Times(1)

// 	mockUserRepo.EXPECT().
// 		Update(gomock.Any(), gomock.Any()).
// 		DoAndReturn(func(db *gorm.DB, u *domain.User) error {
// 			assert.NotEqual(t, "oldhashedpassword", u.Password)
// 			assert.NotEmpty(t, u.Password)
// 			return nil
// 		}).Times(1)

// 	mockAuthRepo.EXPECT().
// 		RevokePasswordResetByToken(gomock.Any(), token).
// 		Return(nil).Times(1)

// 	err := usecase.ResetPassword(context.Background(), token, newPassword)

// 	assert.NoError(t, err)
// }
