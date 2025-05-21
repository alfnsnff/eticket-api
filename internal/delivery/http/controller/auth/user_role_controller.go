package controller

import (
	authmodel "eticket-api/internal/model/auth"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/utils/helper/meta"
	"eticket-api/pkg/utils/helper/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserRoleController struct {
	UserRoleUsecase *authusecase.UserRoleUsecase
}

// NewUserRoleRoleController creates a new UserRoleRoleController instance
func NewUserRoleController(role_usecase *authusecase.UserRoleUsecase) *UserRoleController {
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
	params := meta.GetParams(ctx)
	datas, total, err := urc.UserRoleUsecase.GetAllUserRoles(ctx, params.Limit, params.Offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve user roles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewPaginatedResponse(datas, "User roles retrieved successfully", total, params.Limit, params.Page))
}

func (urc *UserRoleController) GetUserRoleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user role ID", err.Error()))
		return
	}

	data, err := urc.UserRoleUsecase.GetUserRoleByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve user role", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("User role not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "User role retrieved successfully", nil))
}

func (urc *UserRoleController) UpdateUserRole(ctx *gin.Context) {
	request := new(authmodel.UpdateUserRoleRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("User role ID is required", nil))
		return
	}

	if err := urc.UserRoleUsecase.UpdateUserRole(ctx, uint(id), request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update user role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "User role updated successfully", nil))
}

func (urc *UserRoleController) DeleteUserRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user role ID", err.Error()))
		return
	}

	if err := urc.UserRoleUsecase.DeleteUserRole(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete user role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "User role deleted successfully", nil))
}
