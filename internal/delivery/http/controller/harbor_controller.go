package controller

import (
	"eticket-api/internal/model"
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
func NewHarborController(harbor_usecase *usecase.HarborUsecase) *HarborController {
	return &HarborController{HarborUsecase: harbor_usecase}
}

// CreateHarbor handles creating a new harbor
func (h *HarborController) CreateHarbor(ctx *gin.Context) {
	request := new(model.WriteHarborRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	// harbor := dto.ToHarborEntity(request)

	if err := h.HarborUsecase.CreateHarbor(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Harbor created successfully", nil))
}

// GetAllHarbors handles retrieving all harbors
func (h *HarborController) GetAllHarbors(ctx *gin.Context) {
	datas, err := h.HarborUsecase.GetAllHarbors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbors", err.Error()))
		return
	}

	// harborDTOs := dto.ToHarborDTOs(harbors)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Harbors retrieved successfully", nil))
}

// GetHarborByID handles retrieving a harbor by its ID
func (h *HarborController) GetHarborByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	data, err := h.HarborUsecase.GetHarborByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbor", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Harbor not found", nil))
		return
	}

	// harborDTO := dto.ToHarborDTO(harbor)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Harbor retrieved successfully", nil))
}

// UpdateHarbor handles updating an existing harbor
func (h *HarborController) UpdateHarbor(ctx *gin.Context) {
	request := new(model.WriteHarborRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Harbor ID is required", nil))
		return
	}

	// harbor := dto.ToHarborEntity(request)

	if err := h.HarborUsecase.UpdateHarbor(ctx, uint(id), request); err != nil {
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
