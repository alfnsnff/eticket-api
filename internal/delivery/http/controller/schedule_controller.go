package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/schedule"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	Log             logger.Logger
	Validate        validator.Validator
	ScheduleUsecase *schedule.ScheduleUsecase
}

func NewScheduleController(
	log logger.Logger,
	validate validator.Validator,
	schedule_usecase *schedule.ScheduleUsecase,
) *ScheduleController {
	return &ScheduleController{
		Log:             log,
		Validate:        validate,
		ScheduleUsecase: schedule_usecase,
	}

}

func (scc *ScheduleController) CreateSchedule(ctx *gin.Context) {
	request := new(model.WriteScheduleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}

func (scc *ScheduleController) GetAllSchedules(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := scc.ScheduleUsecase.GetAllSchedules(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

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

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Schedules retrieved successfully", total, params.Limit, params.Page))
}

func (scc *ScheduleController) GetAllScheduled(ctx *gin.Context) {
	datas, err := scc.ScheduleUsecase.GetAllScheduled(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Schedules retrieved successfully", nil))
}

func (scc *ScheduleController) GetScheduleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	data, err := scc.ScheduleUsecase.GetScheduleByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedule", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Schedule not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Schedule retrieved successfully", nil))
}

func (scc *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateScheduleRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := scc.Validate.Struct(request); err != nil {
		scc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := scc.ScheduleUsecase.UpdateSchedule(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule updated successfully", nil))
}

func (scc *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	if err := scc.ScheduleUsecase.DeleteSchedule(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule deleted successfully", nil))
}

func (scc *ScheduleController) GetQuotaByScheduleID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	data, err := scc.ScheduleUsecase.GetScheduleAvailability(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedule", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Schedule not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Schedule retrieved successfully", nil))
}

func (scc *ScheduleController) CreateScheduleWithAllocation(ctx *gin.Context) {
	request := new(model.WriteScheduleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := scc.Validate.Struct(request); err != nil {
		scc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := scc.ScheduleUsecase.CreateScheduleWithAllocation(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}
