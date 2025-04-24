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
	PriceUsecase usecase.PriceUsecase
}

// CreateShip handles creating a new Ship
func (h *PriceController) CreatePrice(c *gin.Context) {
	var priceCreate dto.PriceCreate
	if err := c.ShouldBindJSON(&priceCreate); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	price := dto.ToPriceEntity(&priceCreate)

	if err := h.PriceUsecase.CreatePrice(&price); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *PriceController) GetAllPrices(c *gin.Context) {
	prices, err := h.PriceUsecase.GetAllPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	priceDTOs := dto.ToPriceDTOs(prices)
	c.JSON(http.StatusOK, response.NewSuccessResponse(priceDTOs, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *PriceController) GetPriceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	price, err := h.PriceUsecase.GetPriceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if price == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	priceDTO := dto.ToPriceDTO(price)
	c.JSON(http.StatusOK, response.NewSuccessResponse(priceDTO, "Ship retrieved successfully", nil))
}

func (h *PriceController) GetPriceByRouteID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	price, err := h.PriceUsecase.GetPriceByRouteID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if price == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	priceDTO := dto.ToPriceDTOs(price)
	c.JSON(http.StatusOK, response.NewSuccessResponse(priceDTO, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *PriceController) UpdatePrice(c *gin.Context) {
	var priceUpdate dto.PriceCreate

	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&priceUpdate); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	price := dto.ToPriceEntity(&priceUpdate)

	if err := h.PriceUsecase.UpdatePrice(uint(id), &price); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Price updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *PriceController) DeletePrice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.PriceUsecase.DeletePrice(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}

// PriceController creates a new PriceController instance.  Important!
func NewPriceController(priceUsecase usecase.PriceUsecase) *PriceController {
	return &PriceController{PriceUsecase: priceUsecase}
}
