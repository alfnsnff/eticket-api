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

type RoleController struct {
	RoleUsecase *authusecase.RoleUsecase
}

// NewRoleController creates a new RoleController instance
func NewRoleController(role_usecase *authusecase.RoleUsecase) *RoleController {
	return &RoleController{RoleUsecase: role_usecase}
}

func (rc *RoleController) CreateRole(ctx *gin.Context) {
	request := new(authmodel.WriteRoleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := rc.RoleUsecase.CreateRole(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create role", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "role created successfully", nil))
}

func (rc *RoleController) GetAllRoles(ctx *gin.Context) {
	params := meta.GetParams(ctx)
	datas, total, err := rc.RoleUsecase.GetAllRoles(ctx, params.Limit, params.Offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve roles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewPaginatedResponse(datas, "Roles retrieved successfully", total, params.Limit, params.Page))
}

func (rc *RoleController) GetRoleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid role ID", err.Error()))
		return
	}

	data, err := rc.RoleUsecase.GetRoleByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve role", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Role not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Role retrieved successfully", nil))
}

func (rc *RoleController) UpdateRole(ctx *gin.Context) {
	request := new(authmodel.UpdateRoleRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Role ID is required", nil))
		return
	}

	if err := rc.RoleUsecase.UpdateRole(ctx, uint(id), request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Role updated successfully", nil))
}

func (rc *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid role ID", err.Error()))
		return
	}

	if err := rc.RoleUsecase.DeleteRole(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Role deleted successfully", nil))
}
