package v1

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Validate    validator.Validator
	Log         logger.Logger
	RoleUsecase *usecase.RoleUsecase
}

// NewRoleController creates a new RoleController instance
func NewRoleController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	role_usecase *usecase.RoleUsecase,

) {
	roc := &RoleController{
		Log:         log,
		Validate:    validate,
		RoleUsecase: role_usecase,
	}

	router.GET("/roles", roc.GetAllRoles)
	router.GET("/role/:id", roc.GetRoleByID)

	router.POST("/role/create", roc.CreateRole)
	protected.PUT("/role/update/:id", roc.UpdateRole)
	protected.DELETE("/role/:id", roc.DeleteRole)
}

func (rc *RoleController) CreateRole(ctx *gin.Context) {

	request := new(model.WriteRoleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		rc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := rc.Validate.Struct(request); err != nil {
		rc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := rc.RoleUsecase.CreateRole(ctx, request); err != nil {
		rc.Log.WithError(err).Error("failed to create role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create role", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "role created successfully", nil))
}

func (rc *RoleController) GetAllRoles(ctx *gin.Context) {

	params := response.GetParams(ctx)
	datas, total, err := rc.RoleUsecase.ListRoles(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		rc.Log.WithError(err).Error("failed to retrieve roles")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve roles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Roles retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (rc *RoleController) GetRoleByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		rc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse role ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid role ID", err.Error()))
		return
	}

	data, err := rc.RoleUsecase.GetRoleByID(ctx, uint(id))

	if err != nil {
		rc.Log.WithError(err).WithField("id", id).Error("failed to retrieve role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve role", err.Error()))
		return
	}

	if data == nil {
		rc.Log.WithField("id", id).Warn("role not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Role not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Role retrieved successfully", nil))
}

func (rc *RoleController) UpdateRole(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		rc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse role ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing role ID", nil))
		return
	}

	request := new(model.UpdateRoleRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		rc.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := rc.Validate.Struct(request); err != nil {
		rc.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := rc.RoleUsecase.UpdateRole(ctx, request); err != nil {
		rc.Log.WithError(err).WithField("id", id).Error("failed to update role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Role updated successfully", nil))
}

func (rc *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		rc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse role ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid role ID", err.Error()))
		return
	}

	if err := rc.RoleUsecase.DeleteRole(ctx, uint(id)); err != nil {
		rc.Log.WithError(err).WithField("id", id).Error("failed to delete role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Role deleted successfully", nil))
}
