package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManifestController struct {
	ManifestUsecase *usecase.ManifestUsecase
}

// NewManifestController creates a new ManifestController instance.  Important!
func NewManifestController(manifest_usecase *usecase.ManifestUsecase) *ManifestController {
	return &ManifestController{ManifestUsecase: manifest_usecase}
}

// CreateShip handles creating a new Ship
func (h *ManifestController) CreateManifest(ctx *gin.Context) {
	// var shipClassCreate dto.ShipClassCreate
	request := new(model.WriteManifestRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	// shipClass := dto.ToShipClassEntity(&shipClassCreate)

	if err := h.ManifestUsecase.CreateManifest(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *ManifestController) GetAllManifests(ctx *gin.Context) {
	datas, err := h.ManifestUsecase.GetAllManifests(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	// shipClassDTOs := dto.ToShipClassDTOs(shipClasses)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ManifestController) GetManifestByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := h.ManifestUsecase.GetManifestByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	// shipDTO := dto.ToShipClassDTO(shipClass)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ManifestController) GetManifestsByShipID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	capacities, err := h.ManifestUsecase.GetManifestsByShipID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if capacities == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship Class not found", nil))
		return
	}

	// shipClassDTOs := dto.ToShipClassDTOs(shipClasses)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(capacities, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *ManifestController) UpdateManifest(ctx *gin.Context) {
	request := new(model.WriteManifestRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	// shipClass := dto.ToShipClassEntity(&shipClassUpdate)

	if err := h.ManifestUsecase.UpdateManifest(ctx, uint(id), request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *ManifestController) DeleteManifest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.ManifestUsecase.DeleteManifest(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
