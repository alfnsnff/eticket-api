package controller

import (
	authmodel "eticket-api/internal/model/auth"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/utils/helper/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserRoleController struct {
	UserRoleUsecase *authusecase.UserRoleUsecase
}

// NewUserRoleRoleController creates a new UserRoleRoleController instance
func NewUserRoleRoleController(role_usecase *authusecase.UserRoleUsecase) *UserRoleController {
	return &UserRoleController{UserRoleUsecase: role_usecase}
}

func (urc *UserRoleController) CreateUserRole(ctx *gin.Context) {
	request := new(authmodel.WriteUserRoleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := urc.UserRoleUsecase.CreateUserRole(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to assign userRole role", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "UserRole role assigned successfully", nil))
}

func (urc *UserRoleController) GetAllUserRoles(ctx *gin.Context) {

	datas, err := urc.UserRoleUsecase.GetAllUserRoles(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "UserRoles retrieved successfully", nil))
}

func (urc *UserRoleController) GetUserRoleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := urc.UserRoleUsecase.GetUserRoleByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

func (urc *UserRoleController) UpdateUserRole(ctx *gin.Context) {
	request := new(authmodel.UpdateUserRoleRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	if err := urc.UserRoleUsecase.UpdateUserRole(ctx, uint(id), request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

func (urc *UserRoleController) DeleteUserRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := urc.UserRoleUsecase.DeleteUserRole(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
