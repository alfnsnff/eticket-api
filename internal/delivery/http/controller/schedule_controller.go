package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/helper/meta"
	"eticket-api/pkg/utils/helper/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	ScheduleUsecase *usecase.ScheduleUsecase
}

func NewScheduleController(schedule_usecase *usecase.ScheduleUsecase) *ScheduleController {
	return &ScheduleController{ScheduleUsecase: schedule_usecase}
}

func (scc *ScheduleController) CreateSchedule(ctx *gin.Context) {
	request := new(model.WriteScheduleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := scc.ScheduleUsecase.CreateSchedule(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}

func (scc *ScheduleController) GetAllSchedules(ctx *gin.Context) {
	params := meta.GetParams(ctx)
	datas, total, err := scc.ScheduleUsecase.GetAllSchedules(ctx, params.Limit, params.Offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewPaginatedResponse(datas, "Schedules retrieved successfully", total, params.Limit, params.Page))
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
	request := new(model.UpdateScheduleRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Schedule ID is required", nil))
		return
	}

	if err := scc.ScheduleUsecase.UpdateSchedule(ctx, uint(id), request); err != nil {
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

	if err := scc.ScheduleUsecase.CreateScheduleWithAllocation(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}
