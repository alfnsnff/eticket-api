package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	requests "eticket-api/internal/delivery/http/v1/request"
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
	c := &RoleController{
		Log:         log,
		Validate:    validate,
		RoleUsecase: role_usecase,
	}

	router.GET("/roles", c.GetAllRoles)
	router.GET("/role/:id", c.GetRoleByID)

	router.POST("/role/create", c.CreateRole)
	protected.PUT("/role/update/:id", c.UpdateRole)
	protected.DELETE("/role/:id", c.DeleteRole)
}

func (c *RoleController) CreateRole(ctx *gin.Context) {

	request := new(requests.CreateRoleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.RoleUsecase.CreateRole(ctx, requests.RoleFromCreate(request)); err != nil {

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}
		c.Log.WithError(err).Error("failed to create role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create role", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "role created successfully", nil))
}

func (c *RoleController) GetAllRoles(ctx *gin.Context) {

	params := response.GetParams(ctx)
	datas, total, err := c.RoleUsecase.ListRoles(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve roles")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve roles", err.Error()))
		return
	}

	responses := make([]*requests.RoleResponse, len(datas))
	for i, data := range datas {
		responses[i] = requests.RoleToResponse(data)
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		responses,
		"Roles retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (c *RoleController) GetRoleByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse role ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid role ID", err.Error()))
		return
	}

	data, err := c.RoleUsecase.GetRoleByID(ctx, uint(id))

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("role not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("role not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(requests.RoleToResponse(data), "Role retrieved successfully", nil))
}

func (c *RoleController) UpdateRole(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse role ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing role ID", nil))
		return
	}

	request := new(requests.UpdateRoleRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.RoleUsecase.UpdateRole(ctx, requests.RoleFromUpdate(request)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("role not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("role not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("role already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("role already exists", nil))
			return
		}
		c.Log.WithError(err).WithField("id", id).Error("failed to update role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Role updated successfully", nil))
}

func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse role ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid role ID", err.Error()))
		return
	}

	if err := c.RoleUsecase.DeleteRole(ctx, uint(id)); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to delete role")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete role", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Role deleted successfully", nil))
}
