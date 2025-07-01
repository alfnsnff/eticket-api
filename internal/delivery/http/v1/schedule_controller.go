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

type ScheduleController struct {
	Validate        validator.Validator
	Log             logger.Logger
	ScheduleUsecase *usecase.ScheduleUsecase
}

func NewScheduleController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	schedule_usecase *usecase.ScheduleUsecase,

) {
	c := &ScheduleController{
		Log:             log,
		Validate:        validate,
		ScheduleUsecase: schedule_usecase,
	}

	router.GET("/schedules", c.GetAllSchedules)
	router.GET("/schedules/active", c.GetAllScheduled)
	router.GET("/schedule/:id", c.GetScheduleByID)

	protected.POST("/schedule/create", c.CreateSchedule)
	protected.PUT("/schedule/update/:id", c.UpdateSchedule)
	protected.DELETE("/schedule/:id", c.DeleteSchedule)
}

func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {

	request := new(requests.CreateScheduleRequest)

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

	if err := c.ScheduleUsecase.CreateSchedule(ctx, requests.ScheduleFromCreate(request)); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}

		c.Log.WithError(err).Error("failed to create schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}

func (c *ScheduleController) GetAllSchedules(ctx *gin.Context) {

	params := response.GetParams(ctx)
	datas, total, err := c.ScheduleUsecase.ListSchedules(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve schedules")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	responses := make([]*requests.ScheduleResponse, len(datas))
	for i, data := range datas {
		responses[i] = requests.ScheduleToResponse(data)
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		responses,
		"Schedules retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (c *ScheduleController) GetAllScheduled(ctx *gin.Context) {

	datas, err := c.ScheduleUsecase.ListActiveSchedules(ctx)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve active schedules")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	responses := make([]*requests.ScheduleResponse, len(datas))
	for i, data := range datas {
		responses[i] = requests.ScheduleToResponse(data)
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(responses, "Schedules retrieved successfully", nil))
}

func (c *ScheduleController) GetScheduleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	data, err := c.ScheduleUsecase.GetScheduleByID(ctx, uint(id))

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("schedule not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("schedule not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(requests.ScheduleToResponse(data), "Schedule retrieved successfully", nil))
}

func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing schedule ID", nil))
		return
	}

	request := new(requests.UpdateScheduleRequest)
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

	if err := c.ScheduleUsecase.UpdateSchedule(ctx, requests.ScheduleFromUpdate(request)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("schedule not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("schedule not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("schedule already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("schedule already exists", nil))
			return
		}
		c.Log.WithError(err).WithField("id", id).Error("failed to update schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule updated successfully", nil))
}

func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	if err := c.ScheduleUsecase.DeleteSchedule(ctx, uint(id)); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to delete schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule deleted successfully", nil))
}
