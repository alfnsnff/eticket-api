package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model" // Import the response package
	"eticket-api/internal/usecase/fare"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FareController struct {
	Log         logger.Logger
	Validate    validator.Validator
	FareUsecase *fare.FareUsecase
}

func NewFareController(
	log logger.Logger,
	validate validator.Validator,
	Fare_usecase *fare.FareUsecase,
) *FareController {
	return &FareController{
		Log:         log,
		Validate:    validate,
		FareUsecase: Fare_usecase,
	}

}

func (fc *FareController) CreateFare(ctx *gin.Context) {
	request := new(model.WriteFareRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := fc.Validate.Struct(request); err != nil {
		fc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := fc.FareUsecase.CreateFare(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create fare", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Fare created successfully", nil))
}

func (fc *FareController) GetAllFares(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := fc.FareUsecase.GetAllFares(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve fares", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Fares retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Fares retrieved successfully", total, params.Limit, params.Page))
}

func (fc *FareController) GetFareByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid fare ID", err.Error()))
		return
	}

	data, err := fc.FareUsecase.GetFareByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve fare", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Fare not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Fare retrieved successfully", nil))
}

func (fc *FareController) UpdateFare(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateFareRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := fc.Validate.Struct(request); err != nil {
		fc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := fc.FareUsecase.UpdateFare(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update fare", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Fare updated successfully", nil))
}

func (fc *FareController) DeleteFare(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid fare ID", err.Error()))
		return
	}

	if err := fc.FareUsecase.DeleteFare(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete fare", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Fare deleted successfully", nil))
}
