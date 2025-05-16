package controller

import (
	"eticket-api/config"
	authmodel "eticket-api/internal/model/auth"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/utils/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase *authusecase.UserUsecase
}

// NewUserController creates a new UserController instance
func NewUserController(user_usecase *authusecase.UserUsecase) *UserController {
	return &UserController{UserUsecase: user_usecase}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	request := new(authmodel.WriteUserRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := uc.UserUsecase.CreateUser(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create user", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "User created successfully", nil))
}

func (uc *UserController) Login(ctx *gin.Context) {
	request := new(authmodel.UserLoginRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	accessToken, refreshToken, err := uc.UserUsecase.Login(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials", err.Error()))
		return
	}

	// OPTIONAL: Set as HTTP-only secure cookies
	ctx.SetCookie("access_token", accessToken, int(config.AppConfig.Auth.AccessTokenExpiry)*60, "/", "", true, true)
	ctx.SetCookie("refresh_token", refreshToken, int(config.AppConfig.Auth.RefreshTokenExpiry)*24*60*60, "/", "", true, true)

	// OR: Return tokens in JSON (useful for SPA apps)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Login successful", nil))
}
