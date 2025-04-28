package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	ScheduleUsecase *usecase.ScheduleUsecase
}

// NewScheduleController creates a new ScheduleController instance.  Important!
func NewScheduleController(scheduleUsecase *usecase.ScheduleUsecase) *ScheduleController {
	return &ScheduleController{ScheduleUsecase: scheduleUsecase}
}

// CreateSchedule handles creating a new Schedule
func (h *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var scheduleCreate dto.ScheduleCreate
	if err := ctx.ShouldBindJSON(&scheduleCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	schedule := dto.ToScheduleEntity(&scheduleCreate)

	if err := h.ScheduleUsecase.CreateSchedule(ctx, &schedule); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}

// GetAllSchedules handles retrieving all Schedules
func (h *ScheduleController) GetAllSchedules(ctx *gin.Context) {
	schedules, err := h.ScheduleUsecase.GetAllSchedules(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	scheduleDTOs := dto.ToScheduleDTOs(schedules)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(scheduleDTOs, "Schedules retrieved successfully", nil))
}

// GetScheduleByID handles retrieving a Schedule by its ID
func (h *ScheduleController) GetScheduleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	schedule, err := h.ScheduleUsecase.GetScheduleByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedule", err.Error()))
		return
	}

	if schedule == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Schedule not found", nil))
		return
	}

	scheduleDTO := dto.ToScheduleDTO(schedule)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(scheduleDTO, "Schedule retrieved successfully", nil))
}

func (h *ScheduleController) GetPricesWithQuota(ctx *gin.Context) {
	scheduleIDParam := ctx.Param("id")
	scheduleID, err := strconv.ParseUint(scheduleIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	data, err := h.ScheduleUsecase.GetPricesWithQuotaBySchedule(ctx, uint(scheduleID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Schedule retrieved successfully", nil))
}

func (h *ScheduleController) SearchSchedule(ctx *gin.Context) {
	var req dto.ScheduleSearchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.ScheduleUsecase.SearchSchedule(ctx, req)

	scheduleDTO := dto.ToScheduleDTO(schedule)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(scheduleDTO, "Schedule retrieved successfully", nil))
}

// UpdateSchedule handles updating an existing Schedule
func (h *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	var scheduleUpdate dto.ScheduleCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&scheduleUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Schedule ID is required", nil))
		return
	}

	schedule := dto.ToScheduleEntity(&scheduleUpdate)

	if err := h.ScheduleUsecase.UpdateSchedule(ctx, uint(id), &schedule); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule updated successfully", nil))
}

// DeleteSchedule handles deleting a Schedule by its ID
func (h *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	if err := h.ScheduleUsecase.DeleteSchedule(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete schedule", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule deleted successfully", nil))
}
