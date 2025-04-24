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
	ScheduleUsecase usecase.ScheduleUsecase
}

// CreateSchedule handles creating a new Schedule
func (h *ScheduleController) CreateSchedule(c *gin.Context) {
	var scheduleCreate dto.ScheduleCreate
	if err := c.ShouldBindJSON(&scheduleCreate); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	schedule := dto.ToScheduleEntity(&scheduleCreate)

	if err := h.ScheduleUsecase.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create schedule", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Schedule created successfully", nil))
}

// GetAllSchedules handles retrieving all Schedules
func (h *ScheduleController) GetAllSchedules(c *gin.Context) {
	schedules, err := h.ScheduleUsecase.GetAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedules", err.Error()))
		return
	}

	scheduleDTOs := dto.ToScheduleDTOs(schedules)
	c.JSON(http.StatusOK, response.NewSuccessResponse(scheduleDTOs, "Schedules retrieved successfully", nil))
}

// GetScheduleByID handles retrieving a Schedule by its ID
func (h *ScheduleController) GetScheduleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	schedule, err := h.ScheduleUsecase.GetScheduleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve schedule", err.Error()))
		return
	}

	if schedule == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Schedule not found", nil))
		return
	}

	scheduleDTO := dto.ToScheduleDTO(schedule)
	c.JSON(http.StatusOK, response.NewSuccessResponse(scheduleDTO, "Schedule retrieved successfully", nil))
}

func (h *ScheduleController) GetPricesWithQuota(c *gin.Context) {
	scheduleIDParam := c.Param("id")
	scheduleID, err := strconv.ParseUint(scheduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	data, err := h.ScheduleUsecase.GetPricesWithQuotaBySchedule(uint(scheduleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(data, "Schedule retrieved successfully", nil))
}

func (h *ScheduleController) SearchSchedule(c *gin.Context) {
	var req dto.ScheduleSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.ScheduleUsecase.SearchSchedule(req)

	scheduleDTO := dto.ToScheduleDTO(schedule)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(scheduleDTO, "Schedule retrieved successfully", nil))
}

// UpdateSchedule handles updating an existing Schedule
func (h *ScheduleController) UpdateSchedule(c *gin.Context) {
	var scheduleUpdate dto.ScheduleCreate

	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&scheduleUpdate); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Schedule ID is required", nil))
		return
	}

	schedule := dto.ToScheduleEntity(&scheduleUpdate)

	if err := h.ScheduleUsecase.UpdateSchedule(uint(id), &schedule); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update schedule", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule updated successfully", nil))
}

// DeleteSchedule handles deleting a Schedule by its ID
func (h *ScheduleController) DeleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	if err := h.ScheduleUsecase.DeleteSchedule(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete schedule", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule deleted successfully", nil))
}

// NewScheduleController creates a new ScheduleController instance.  Important!
func NewScheduleController(scheduleUsecase usecase.ScheduleUsecase) *ScheduleController {
	return &ScheduleController{ScheduleUsecase: scheduleUsecase}
}
