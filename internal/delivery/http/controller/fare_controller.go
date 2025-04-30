package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FareController struct {
	FareUsecase *usecase.FareUsecase
}

// FareController creates a new FareController instance.  Important!
func NewFareController(Fare_usecase *usecase.FareUsecase) *FareController {
	return &FareController{FareUsecase: Fare_usecase}
}

// CreateShip handles creating a new Ship
func (h *FareController) CreateFare(ctx *gin.Context) {

	request := new(model.WriteFareRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Fare := dto.ToFareEntity(&FareCreate)

	if err := h.FareUsecase.CreateFare(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

// GetAllShips handles retrieving all Ships
func (h *FareController) GetAllFares(ctx *gin.Context) {
	datas, err := h.FareUsecase.GetAllFares(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	// FareDTOs := dto.ToFareDTOs(Fares)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Ships retrieved successfully", nil))
}

// GetShipByID handles retrieving a Ship by its ID
func (h *FareController) GetFareByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := h.FareUsecase.GetFareByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	// FareDTO := dto.ToFareDTO(Fare)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

func (h *FareController) GetFareByRouteID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := h.FareUsecase.GetFareByRouteID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	// FareDTO := dto.ToFareDTOs(Fare)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

// UpdateShip handles updating an existing Ship
func (h *FareController) UpdateFare(ctx *gin.Context) {
	request := new(model.WriteFareRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}

	// Fare := dto.ToFareEntity(&FareUpdate)

	if err := h.FareUsecase.UpdateFare(ctx, uint(id), request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Fare updated successfully", nil))
}

// DeleteShip handles deleting a Ship by its ID
func (h *FareController) DeleteFare(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := h.FareUsecase.DeleteFare(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
