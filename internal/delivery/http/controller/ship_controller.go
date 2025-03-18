package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipController struct {
	ShipUsecase usecase.ShipUsecase
}

// CreateShip handles creating a new Ship
func (h *ShipController) CreateShip(c *gin.Context) {
	var ship entities.Ship
	if err := c.ShouldBindJSON(&ship); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := h.ShipUsecase.CreateShip(&ship); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *ShipController) GetAllShips(c *gin.Context) {
	ships, err := h.ShipUsecase.GetAllShips()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	shipDTOs := dto.ToShipDTOs(ships)
	c.JSON(http.StatusOK, response.NewSuccessResponse(shipDTOs, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *ShipController) GetShipByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	ship, err := h.ShipUsecase.GetShipByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if ship == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	shipDTO := dto.ToShipDTO(ship)
	c.JSON(http.StatusOK, response.NewSuccessResponse(shipDTO, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *ShipController) UpdateShip(c *gin.Context) {
	var ship entities.Ship
	if err := c.ShouldBindJSON(&ship); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if ship.ID == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	if err := h.ShipUsecase.UpdateShip(&ship); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *ShipController) DeleteShip(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.ShipUsecase.DeleteShip(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}

// NewShipController creates a new ShipController instance.  Important!
func NewShipController(shipUsecase usecase.ShipUsecase) *ShipController {
	return &ShipController{ShipUsecase: shipUsecase}
}
