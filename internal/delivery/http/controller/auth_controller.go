package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Validate     validator.Validator
	Log          logger.Logger
	TokenManager *token.JWT
	AuthUsecase  *auth.AuthUsecase
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

// NewUserRoleController creates a new UserRoleController instance
func NewAuthController(
	router *gin.Engine,
	log logger.Logger,
	validate validator.Validator,
	auth_usecase *auth.AuthUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	ac := &AuthController{
		Log:          log,
		Validate:     validate,
		AuthUsecase:  auth_usecase,
		Authenticate: authtenticate,
		Authorized:   authorized,
	}

	public := router.Group("/api/v1") // No middleware
	public.GET("/auth/me", ac.Me)
	public.POST("/auth/login", ac.Login)
	public.POST("/auth/refresh", ac.RefreshToken)
	public.POST("/auth/forget-password", ac.ForgetPassword)

	protected := router.Group("/api/v1")
	protected.Use(ac.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/auth/logout", ac.Logout)
}

func (auc *AuthController) Login(ctx *gin.Context) {
	request := new(model.WriteLoginRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := auc.Validate.Struct(request); err != nil {
		auc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, accessToken, refreshToken, err := auc.AuthUsecase.Login(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials", err.Error()))
		return
	}

	// OPTIONAL: Set as HTTP-only secure cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", accessToken, 15*60, "/", "", true, true)      // 15 minutes
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", true, true) // 1 day

	// OR: Return tokens in JSON (useful for SPA apps)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Login successful", nil))
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

func (auc *AuthController) ForgetPassword(ctx *gin.Context) {
	request := new(model.WriteForgetPasswordRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := auc.Validate.Struct(request); err != nil {
		auc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := auc.AuthUsecase.RequestPasswordReset(ctx, request.Email); err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Reset password failed", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "We will send reset password email if it matched to our system", nil))
}

func (auc *AuthController) Me(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing access token", err.Error()))
		return
	}

	user, err := auc.AuthUsecase.Me(ctx, accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Unauthorized", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(user, "User info retrieved successfully", nil))
}
