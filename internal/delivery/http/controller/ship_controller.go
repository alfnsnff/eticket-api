package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipController struct {
	ShipUsecase *usecase.ShipUsecase
}

// NewShipController creates a new ShipController instance.  Important!
func NewShipController(shipUsecase *usecase.ShipUsecase) *ShipController {
	return &ShipController{ShipUsecase: shipUsecase}
}

// CreateShip handles creating a new Ship
func (h *ShipController) CreateShip(ctx *gin.Context) {
	var shipCreate dto.ShipCreate
	if err := ctx.ShouldBindJSON(&shipCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	ship := dto.ToShipEntity(&shipCreate)

	if err := h.ShipUsecase.CreateShip(ctx, &ship); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *ShipController) GetAllShips(ctx *gin.Context) {
	ships, err := h.ShipUsecase.GetAllShips(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	shipDTOs := dto.ToShipDTOs(ships)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(shipDTOs, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ShipController) GetShipByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	ship, err := h.ShipUsecase.GetShipByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if ship == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	shipDTO := dto.ToShipDTO(ship)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(shipDTO, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *ShipController) UpdateShip(ctx *gin.Context) {
	var shipUpdate dto.ShipCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&shipUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	ship := dto.ToShipEntity(&shipUpdate)

	if err := h.ShipUsecase.UpdateShip(ctx, uint(id), &ship); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *ShipController) DeleteShip(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.ShipUsecase.DeleteShip(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
