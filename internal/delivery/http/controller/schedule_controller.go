package controller

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
	scc := &ScheduleController{
		Log:             log,
		Validate:        validate,
		ScheduleUsecase: schedule_usecase,
	}

	router.GET("/schedules", scc.GetAllSchedules)
	router.GET("/schedules/active", scc.GetAllScheduled)
	router.GET("/schedule/:id", scc.GetScheduleByID)

	protected.POST("/schedule/create", scc.CreateSchedule)
	protected.PUT("/schedule/update/:id", scc.UpdateSchedule)
	protected.DELETE("/schedule/:id", scc.DeleteSchedule)
}

func (scc *ScheduleController) CreateSchedule(ctx *gin.Context) {

	request := new(model.WriteScheduleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		scc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := scc.Validate.Struct(request); err != nil {
		scc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := scc.ScheduleUsecase.CreateSchedule(ctx, request); err != nil {
		scc.Log.WithError(err).Error("failed to create schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}

func (scc *ScheduleController) GetAllSchedules(ctx *gin.Context) {

	params := response.GetParams(ctx)
	datas, total, err := scc.ScheduleUsecase.ListSchedules(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		scc.Log.WithError(err).Error("failed to retrieve schedules")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	scc.Log.WithField("count", total).Info("Schedules retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Schedules retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (scc *ScheduleController) GetAllScheduled(ctx *gin.Context) {

	datas, err := scc.ScheduleUsecase.ListActiveSchedules(ctx)

	if err != nil {
		scc.Log.WithError(err).Error("failed to retrieve active schedules")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	scc.Log.WithField("count", len(datas)).Info("Active schedules retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Schedules retrieved successfully", nil))
}

func (scc *ScheduleController) GetScheduleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		scc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	data, err := scc.ScheduleUsecase.GetScheduleByID(ctx, uint(id))

	if err != nil {
		scc.Log.WithError(err).WithField("id", id).Error("failed to retrieve schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedule", err.Error()))
		return
	}

	if data == nil {
		scc.Log.WithField("id", id).Warn("schedule not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Schedule not found", nil))
		return
	}

	scc.Log.WithField("id", id).Info("Schedule retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Schedule retrieved successfully", nil))
}

func (scc *ScheduleController) UpdateSchedule(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		scc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing schedule ID", nil))
		return
	}

	request := new(model.UpdateScheduleRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		scc.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := scc.Validate.Struct(request); err != nil {
		scc.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := scc.ScheduleUsecase.UpdateSchedule(ctx, request); err != nil {
		scc.Log.WithError(err).WithField("id", id).Error("failed to update schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update schedule", err.Error()))
		return
	}

	scc.Log.WithField("id", id).Info("Schedule updated successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule updated successfully", nil))
}

func (scc *ScheduleController) DeleteSchedule(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		scc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	if err := scc.ScheduleUsecase.DeleteSchedule(ctx, uint(id)); err != nil {
		scc.Log.WithError(err).WithField("id", id).Error("failed to delete schedule")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete schedule", err.Error()))
		return
	}

	scc.Log.WithField("id", id).Info("Schedule deleted successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule deleted successfully", nil))
}
