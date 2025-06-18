package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/route"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RouteController struct {
	RouteUsecase *route.RouteUsecase
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewRouteController(
	g *gin.Engine, route_usecase *route.RouteUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	rc := &RouteController{
		RouteUsecase: route_usecase,
		Authenticate: authtenticate,
		Authorized:   authorized}

	public := g.Group("") // No middleware
	public.GET("/routes", rc.GetAllRoutes)
	public.GET("/route/:id", rc.GetRouteByID)

	protected := g.Group("")
	protected.Use(rc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/route/create", rc.CreateRoute)
	protected.PUT("/route//update:id", rc.UpdateRoute)
	protected.DELETE("/route/:id", rc.DeleteRoute)
}

func (rc *RouteController) CreateRoute(ctx *gin.Context) {
	request := new(model.WriteRouteRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
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
	request := new(model.UpdateRouteRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Route ID is required", nil))
		return
	}
	request.ID = uint(id)
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
