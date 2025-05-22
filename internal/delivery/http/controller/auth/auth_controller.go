package controller

import (
	"eticket-api/config"
	authmodel "eticket-api/internal/model/auth"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	Cfg          *config.Config
	TokenManager *jwt.TokenManager
	AuthUsecase  *authusecase.AuthUsecase
}

// NewUserRoleController creates a new UserRoleController instance
func NewAuthController(
	cfg *config.Config,
	tm *jwt.TokenManager,
	auth_usecase *authusecase.AuthUsecase,
) *AuthController {
	return &AuthController{
		Cfg:          cfg,
		TokenManager: tm,
		AuthUsecase:  auth_usecase,
	}
}

func (uc *AuthController) Login(ctx *gin.Context) {
	request := new(authmodel.UserLoginRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	accessToken, refreshToken, err := uc.AuthUsecase.Login(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials", err.Error()))
		return
	}

	// OPTIONAL: Set as HTTP-only secure cookies
	ctx.SetCookie("access_token", accessToken, int(uc.Cfg.Auth.AccessTokenExpiry), "/", "", true, true)
	ctx.SetCookie("refresh_token", refreshToken, int(uc.Cfg.Auth.RefreshTokenExpiry), "/", "", true, true)

	// OR: Return tokens in JSON (useful for SPA apps)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Login successful", nil))
}

func (uc *AuthController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing refresh token", err.Error()))
		return
	}

	// Validate token and get claims
	claims, err := uc.TokenManager.ValidateToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid refresh token", err.Error()))
		return
	}

	// Parse token ID (jti)
	tokenID, err := uuid.Parse(claims.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid token ID", err.Error()))
		return
	}

	// Revoke the token in DB
	err = uc.AuthUsecase.RevokeRefreshToken(ctx, tokenID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to revoke token", err.Error()))
		return
	}

	// Clear cookies
	ctx.SetCookie("access_token", "", -1, "/", "", true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Logout successful", nil))
}
