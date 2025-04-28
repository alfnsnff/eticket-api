package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HarborController struct {
	HarborUsecase *usecase.HarborUsecase
}

// NewHarborController creates a new HarborController instance.
func NewHarborController(harborUsecase *usecase.HarborUsecase) *HarborController {
	return &HarborController{HarborUsecase: harborUsecase}
}

// CreateHarbor handles creating a new harbor
func (h *HarborController) CreateHarbor(ctx *gin.Context) {
	var harborCreate dto.HarborCreate
	if err := ctx.ShouldBindJSON(&harborCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	harbor := dto.ToHarborEntity(&harborCreate)

	if err := h.HarborUsecase.CreateHarbor(ctx, &harbor); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Harbor created successfully", nil))
}

// GetAllHarbors handles retrieving all harbors
func (h *HarborController) GetAllHarbors(ctx *gin.Context) {
	harbors, err := h.HarborUsecase.GetAllHarbors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbors", err.Error()))
		return
	}

	harborDTOs := dto.ToHarborDTOs(harbors)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(harborDTOs, "Harbors retrieved successfully", nil))
}

// GetHarborByID handles retrieving a harbor by its ID
func (h *HarborController) GetHarborByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	harbor, err := h.HarborUsecase.GetHarborByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbor", err.Error()))
		return
	}

	if harbor == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Harbor not found", nil))
		return
	}

	harborDTO := dto.ToHarborDTO(harbor)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(harborDTO, "Harbor retrieved successfully", nil))
}

// UpdateHarbor handles updating an existing harbor
func (h *HarborController) UpdateHarbor(ctx *gin.Context) {
	var harborUpdate dto.HarborCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&harborUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Harbor ID is required", nil))
		return
	}

	harbor := dto.ToHarborEntity(&harborUpdate)

	if err := h.HarborUsecase.UpdateHarbor(ctx, uint(id), &harbor); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor updated successfully", nil))
}

// DeleteHarbor handles deleting a harbor by its ID
func (h *HarborController) DeleteHarbor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	if err := h.HarborUsecase.DeleteHarbor(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor deleted successfully", nil))
}
