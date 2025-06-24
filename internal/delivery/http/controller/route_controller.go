package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/route"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RouteController struct {
	Log          logger.Logger
	Validate     validator.Validator
	RouteUsecase *route.RouteUsecase
}

func NewRouteController(
	log logger.Logger,
	validate validator.Validator,
	route_usecase *route.RouteUsecase,
) *RouteController {
	return &RouteController{
		Log:          log,
		Validate:     validate,
		RouteUsecase: route_usecase,
	}

}

func (rc *RouteController) CreateRoute(ctx *gin.Context) {
	request := new(model.WriteRouteRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := rc.Validate.Struct(request); err != nil {
		rc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := rc.RouteUsecase.CreateRoute(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create route", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Route created successfully", nil))
}

func (rc *RouteController) GetAllRoutes(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := rc.RouteUsecase.GetAllRoutes(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve routes", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Routes retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Routes retrieved successfully", total, params.Limit, params.Page))
}

func (rc *RouteController) GetRouteByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	data, err := rc.RouteUsecase.GetRouteByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve route", err.Error())) // More specific error
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Route not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Route retrieved successfully", nil))
}

func (rc *RouteController) UpdateRoute(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateRouteRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := rc.Validate.Struct(request); err != nil {
		rc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := rc.RouteUsecase.UpdateRoute(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update route", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Route updated successfully", nil))
}

func (rc *RouteController) DeleteRoute(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	if err := rc.RouteUsecase.DeleteRoute(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete route", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Route deleted successfully", nil))
}
