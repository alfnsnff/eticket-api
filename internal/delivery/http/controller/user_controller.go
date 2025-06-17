package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase *user.UserUsecase
}

// NewUserController creates a new UserController instance
func NewUserController(user_usecase *user.UserUsecase) *UserController {
	return &UserController{UserUsecase: user_usecase}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	request := new(model.WriteUserRequest)

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

func (uc *UserController) GetAllUsers(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := uc.UserUsecase.GetAllUsers(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve users", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Fares retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "User roles retrieved successfully", total, params.Limit, params.Page))
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
		return
	}

	data, err := uc.UserUsecase.GetUserByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve user", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("User not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "User retrieved successfully", nil))
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	request := new(model.UpdateUserRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("User ID is required", nil))
		return
	}
	request.ID = uint(id)
	if err := uc.UserUsecase.UpdateUser(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update user", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "User updated successfully", nil))
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
		return
	}

	if err := uc.UserUsecase.DeleteUser(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete user", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "User deleted successfully", nil))
}
