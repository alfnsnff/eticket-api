package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model" // Import the response package
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipController struct {
	Validate    validator.Validator
	Log         logger.Logger
	ShipUsecase *usecase.ShipUsecase
}

func NewShipController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	ship_usecase *usecase.ShipUsecase,

) {
	shc := &ShipController{
		Log:         log,
		Validate:    validate,
		ShipUsecase: ship_usecase,
	}

	router.GET("/ships", shc.GetAllShips)
	router.GET("/ship/:id", shc.GetShipByID)

	protected.POST("/ship/create", shc.CreateShip)
	protected.PUT("/ship/update/:id", shc.UpdateShip)
	protected.DELETE("/ship/:id", shc.DeleteShip)
}

func (shc *ShipController) CreateShip(ctx *gin.Context) {

	request := new(model.WriteShipRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		shc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := shc.Validate.Struct(request); err != nil {
		shc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := shc.ShipUsecase.CreateShip(ctx, request); err != nil {
		shc.Log.WithError(err).Error("failed to create ship")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

func (shc *ShipController) GetAllShips(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := shc.ShipUsecase.ListShips(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		shc.Log.WithError(err).Error("failed to retrieve ships")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	shc.Log.WithField("count", total).Info("Ships retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Ships retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Ships retrieved successfully", total, params.Limit, params.Page))
}

func (shc *ShipController) GetShipByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		shc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid ship ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := shc.ShipUsecase.GetShipByID(ctx, uint(id))

	if err != nil {
		shc.Log.WithError(err).WithField("id", id).Error("failed to retrieve ship")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		shc.Log.WithField("id", id).Warn("ship not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	shc.Log.WithField("id", id).Info("Ship retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

func (shc *ShipController) UpdateShip(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		shc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid or missing ship ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateShipRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		shc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := shc.Validate.Struct(request); err != nil {
		shc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := shc.ShipUsecase.UpdateShip(ctx, request); err != nil {
		shc.Log.WithError(err).WithField("id", id).Error("failed to update ship")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	shc.Log.WithField("id", id).Info("Ship updated successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

func (shc *ShipController) DeleteShip(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		shc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid ship ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := shc.ShipUsecase.DeleteShip(ctx, uint(id)); err != nil {
		shc.Log.WithError(err).WithField("id", id).Error("failed to delete ship")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	shc.Log.WithField("id", id).Info("Ship deleted successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
