package controller

import (
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	TokenManager *jwt.TokenUtil
	AuthUsecase  *auth.AuthUsecase
}

// NewUserRoleController creates a new UserRoleController instance
func NewAuthController(
	auth_usecase *auth.AuthUsecase,
) *AuthController {
	return &AuthController{
		AuthUsecase: auth_usecase,
	}
}

func (auc *AuthController) Login(ctx *gin.Context) {
	request := new(model.UserLoginRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	accessToken, refreshToken, err := auc.AuthUsecase.Login(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials", err.Error()))
		return
	}

	// OPTIONAL: Set as HTTP-only secure cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", accessToken, 15*60, "/", "", true, true)      // 15 minutes
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", true, true) // 1 day

	// OR: Return tokens in JSON (useful for SPA apps)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Login successful", nil))
}

func (auc *AuthController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing refresh token", err.Error()))
		return
	}

	// Revoke the token in DB
	err = auc.AuthUsecase.Logout(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to revoke token", err.Error()))
		return
	}

	// Clear cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", "", -1, "/", "", true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Logout successful", nil))
}

func (auc *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing refresh token", err.Error()))
		return
	}

	newAccessToken, err := auc.AuthUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid session", err.Error()))
		return
	}

	// Set new access token cookie
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", newAccessToken, 15*60, "/", "", true, true)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Token refreshed successfully", nil))
}
