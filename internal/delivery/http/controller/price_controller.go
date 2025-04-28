package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PriceController struct {
	PriceUsecase *usecase.PriceUsecase
}

// PriceController creates a new PriceController instance.  Important!
func NewPriceController(priceUsecase *usecase.PriceUsecase) *PriceController {
	return &PriceController{PriceUsecase: priceUsecase}
}

// CreateShip handles creating a new Ship
func (h *PriceController) CreatePrice(ctx *gin.Context) {
	var priceCreate dto.PriceCreate
	if err := ctx.ShouldBindJSON(&priceCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	price := dto.ToPriceEntity(&priceCreate)

	if err := h.PriceUsecase.CreatePrice(ctx, &price); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *PriceController) GetAllPrices(ctx *gin.Context) {
	prices, err := h.PriceUsecase.GetAllPrices(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	priceDTOs := dto.ToPriceDTOs(prices)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(priceDTOs, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *PriceController) GetPriceByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	price, err := h.PriceUsecase.GetPriceByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if price == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	priceDTO := dto.ToPriceDTO(price)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(priceDTO, "Ship retrieved successfully", nil))
}

func (h *PriceController) GetPriceByRouteID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	price, err := h.PriceUsecase.GetPriceByRouteID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if price == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	priceDTO := dto.ToPriceDTOs(price)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(priceDTO, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *PriceController) UpdatePrice(ctx *gin.Context) {
	var priceUpdate dto.PriceCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&priceUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	price := dto.ToPriceEntity(&priceUpdate)

	if err := h.PriceUsecase.UpdatePrice(ctx, uint(id), &price); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Price updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *PriceController) DeletePrice(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.PriceUsecase.DeletePrice(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
