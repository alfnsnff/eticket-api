package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HarborController struct {
	HarborUsecase usecase.HarborUsecase
}

// CreateHarbor handles creating a new harbor
func (h *HarborController) CreateHarbor(c *gin.Context) {
	var harbor entities.Harbor
	if err := c.ShouldBindJSON(&harbor); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := h.HarborUsecase.CreateHarbor(&harbor); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create harbor", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Harbor created successfully", nil))
}

// GetAllHarbors handles retrieving all harbors
func (h *HarborController) GetAllHarbors(c *gin.Context) {
	harbors, err := h.HarborUsecase.GetAllHarbors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbors", err.Error()))
		return
	}

	harborDTOs := dto.ToHarborDTOs(harbors)
	c.JSON(http.StatusOK, response.NewSuccessResponse(harborDTOs, "Harbors retrieved successfully", nil))
}

// GetHarborByID handles retrieving a harbor by its ID
func (h *HarborController) GetHarborByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	harbor, err := h.HarborUsecase.GetHarborByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbor", err.Error()))
		return
	}

	if harbor == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Harbor not found", nil))
		return
	}

	harborDTO := dto.ToHarborDTO(harbor)
	c.JSON(http.StatusOK, response.NewSuccessResponse(harborDTO, "Harbor retrieved successfully", nil))
}

// UpdateHarbor handles updating an existing harbor
func (h *HarborController) UpdateHarbor(c *gin.Context) {
	var harbor entities.Harbor
	if err := c.ShouldBindJSON(&harbor); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if harbor.ID == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Harbor ID is required", nil))
		return
	}
	if err := h.HarborUsecase.UpdateHarbor(&harbor); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update harbor", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor updated successfully", nil))
}

// DeleteHarbor handles deleting a harbor by its ID
func (h *HarborController) DeleteHarbor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	if err := h.HarborUsecase.DeleteHarbor(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete harbor", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor deleted successfully", nil))
}

// NewHarborController creates a new HarborController instance.
func NewHarborController(harborUsecase usecase.HarborUsecase) *HarborController {
	return &HarborController{HarborUsecase: harborUsecase}
}
