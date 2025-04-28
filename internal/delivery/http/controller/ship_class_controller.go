package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipClassController struct {
	ShipClassUsecase *usecase.ShipClassUsecase
}

// NewShipClassController creates a new ShipClassController instance.  Important!
func NewShipClassController(shipClassUsecase *usecase.ShipClassUsecase) *ShipClassController {
	return &ShipClassController{ShipClassUsecase: shipClassUsecase}
}

// CreateShip handles creating a new Ship
func (h *ShipClassController) CreateShipClass(ctx *gin.Context) {
	var shipClassCreate dto.ShipClassCreate
	if err := ctx.ShouldBindJSON(&shipClassCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	shipClass := dto.ToShipClassEntity(&shipClassCreate)

	if err := h.ShipClassUsecase.CreateShipClass(ctx, &shipClass); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *ShipClassController) GetAllShipClasses(ctx *gin.Context) {
	shipClasses, err := h.ShipClassUsecase.GetAllShipClasses(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	shipClassDTOs := dto.ToShipClassDTOs(shipClasses)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(shipClassDTOs, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ShipClassController) GetShipClassByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	shipClass, err := h.ShipClassUsecase.GetShipClassByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if shipClass == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	shipDTO := dto.ToShipClassDTO(shipClass)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(shipDTO, "Ship retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ShipClassController) GetShipClassByShipID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	shipClasses, err := h.ShipClassUsecase.GetShipClassByShipID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if shipClasses == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship Class not found", nil))
		return
	}

	shipClassDTOs := dto.ToShipClassDTOs(shipClasses)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(shipClassDTOs, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *ShipClassController) UpdateShipClass(ctx *gin.Context) {
	var shipClassUpdate dto.ShipClassCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&shipClassUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	shipClass := dto.ToShipClassEntity(&shipClassUpdate)

	if err := h.ShipClassUsecase.UpdateShipClass(ctx, uint(id), &shipClass); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *ShipClassController) DeleteShipClass(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.ShipClassUsecase.DeleteShipClass(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
