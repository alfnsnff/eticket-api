package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	requests "eticket-api/internal/delivery/http/v1/request"
	"eticket-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Validate     validator.Validator
	Log          logger.Logger
	TokenManager *token.JWT
	AuthUsecase  *usecase.AuthUsecase
}

// NewUserRoleController creates a new UserRoleController instance
func NewAuthController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	auth_usecase *usecase.AuthUsecase,

) {
	c := &AuthController{
		Log:         log,
		Validate:    validate,
		AuthUsecase: auth_usecase,
	}

	router.GET("/auth/me", c.Me)
	router.POST("/auth/login", c.Login)
	router.POST("/auth/refresh", c.RefreshToken)
	// router.POST("/auth/forget-password", c.ForgetPassword)

	protected.POST("/auth/logout", c.Logout)
}

func (c *AuthController) Login(ctx *gin.Context) {
	request := new(requests.LoginRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON login request")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate login request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, accessToken, refreshToken, err := c.AuthUsecase.Login(ctx, requests.LoginFromRequest(request))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithError(err).Warn("Invalid credentials")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Invalid credentials", nil))
			return
		}
		c.Log.WithError(err).Error("failed to authenticate user")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Invalid credentials", err.Error()))
		return
	}

	// OPTIONAL: Set as HTTP-only secure cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", accessToken, 15*60, "/", "", true, true)      // 15 minutes
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", true, true) // 1 day

	// OR: Return tokens in JSON (useful for SPA apps)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Login successful", nil))
}

func (c *AuthController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		c.Log.WithError(err).Error("missing refresh token in logout request")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing refresh token", err.Error()))
		return
	}

	// Revoke the token in DB
	err = c.AuthUsecase.Logout(ctx, refreshToken)
	if err != nil {
		c.Log.WithError(err).Error("failed to revoke refresh token")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Failed to revoke token", err.Error()))
		return
	}

	// Clear cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", "", -1, "/", "", true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Logout successful", nil))
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	c.Log.Info("Processing token refresh request")
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		c.Log.WithError(err).Error("missing refresh token in refresh request")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing refresh token", err.Error()))
		return
	}

	newAccessToken, err := c.AuthUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.Log.WithError(err).Error("failed to refresh token")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid session", err.Error()))
		return
	}

	// Set new access token cookie
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", newAccessToken, 15*60, "/", "", true, true)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Token refreshed successfully", nil))
}

// func (c *AuthController) ForgetPassword(ctx *gin.Context) {
// 	c.Log.Info("Processing forget password request")
// 	request := new(model.WriteForgetPasswordRequest)

// 	if err := ctx.ShouldBindJSON(request); err != nil {
// 		c.Log.WithError(err).Error("failed to bind JSON forget password request")
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 		return
// 	}

// 	if err := c.Validate.Struct(request); err != nil {
// 		c.Log.WithError(err).Error("failed to validate forget password request body")
// 		errors := validator.ParseErrors(err)
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
// 		return
// 	}

// 	if err := c.AuthUsecase.RequestPasswordReset(ctx, request.Email); err != nil {
// 		c.Log.WithError(err).WithField("email", request.Email).Error("failed to process password reset request")
// 		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Reset password failed", err.Error()))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "We will send reset password email if it matched to our system", nil))
// }

func (c *AuthController) Me(ctx *gin.Context) {
	c.Log.Info("Retrieving user profile information")
	accessToken, err := ctx.Cookie("access_token")
	if err != nil {
		c.Log.WithError(err).Error("missing access token in profile request")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing access token", err.Error()))
		return
	}

	user, err := c.AuthUsecase.Me(ctx, accessToken)
	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve user profile")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Unauthorized", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(user, "User info retrieved successfully", nil))
}
