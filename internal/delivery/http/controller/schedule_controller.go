package controller

import (
	"eticket-api/internal/domain"
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
	var schedule domain.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

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

	c.JSON(http.StatusOK, response.NewSuccessResponse(schedules, "Schedules retrieved successfully", nil))
}

// GetScheduleByID handles retrieving a Schedule by its ID
func (h *ScheduleController) GetScheduleByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
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

	c.JSON(http.StatusOK, response.NewSuccessResponse(schedule, "Schedule retrieved successfully", nil))
}

// UpdateSchedule handles updating an existing Schedule
func (h *ScheduleController) UpdateSchedule(c *gin.Context) {
	var schedule domain.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if schedule.ID == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Schedule ID is required", nil))
		return
	}

	if err := h.ScheduleUsecase.UpdateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update schedule", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Schedule updated successfully", nil))
}

// DeleteSchedule handles deleting a Schedule by its ID
func (h *ScheduleController) DeleteSchedule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
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
