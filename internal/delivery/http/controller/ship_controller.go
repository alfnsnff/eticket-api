package controller

import (
	"eticket-api/internal/model"
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
func NewShipController(ship_usecase *usecase.ShipUsecase) *ShipController {
	return &ShipController{ShipUsecase: ship_usecase}
}

// CreateShip handles creating a new Ship
func (h *ShipController) CreateShip(ctx *gin.Context) {
	request := new(model.WriteShipRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	// ship := dto.ToShipEntity(&shipCreate)

	if err := h.ShipUsecase.CreateShip(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *ShipController) GetAllShips(ctx *gin.Context) {
	datas, err := h.ShipUsecase.GetAllShips(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	// shipDTOs := dto.ToShipDTOs(ships)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ShipController) GetShipByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := h.ShipUsecase.GetShipByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	// shipDTO := dto.ToShipDTO(ship)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *ShipController) UpdateShip(ctx *gin.Context) {
	request := new(model.WriteShipRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	// ship := dto.ToShipEntity(&shipUpdate)

	if err := h.ShipUsecase.UpdateShip(ctx, uint(id), request); err != nil {
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
