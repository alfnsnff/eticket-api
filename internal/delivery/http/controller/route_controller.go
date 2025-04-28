package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RouteController struct {
	RouteUsecase *usecase.RouteUsecase
}

// NewRouteController creates a new RouteController instance.
func NewRouteController(routeUsecase *usecase.RouteUsecase) *RouteController {
	return &RouteController{RouteUsecase: routeUsecase}
}

// CreateRoute creates a new route
func (h *RouteController) CreateRoute(ctx *gin.Context) {
	var routeCreate dto.RouteCreate
	if err := ctx.ShouldBindJSON(&routeCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	route := dto.ToRouteEntity(&routeCreate)

	if err := h.RouteUsecase.CreateRoute(ctx, &route); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create route", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Route created successfully", nil))
}

// GetAllRoutes retrieves all routes
func (h *RouteController) GetAllRoutes(ctx *gin.Context) {
	routes, err := h.RouteUsecase.GetAllRoutes(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve routes", err.Error()))
		return
	}

	routeDTOs := dto.ToRouteDTOs(routes)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(routeDTOs, "Routes retrieved successfully", nil))
}

// GetRouteByID retrieves a single route by ID
func (h *RouteController) GetRouteByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	route, err := h.RouteUsecase.GetRouteByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve route", err.Error())) // More specific error
		return
	}
	if route == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Route not found", nil))
		return
	}

	routeDTO := dto.ToRouteDTO(route)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(routeDTO, "Route retrieved successfully", nil))
}

// UpdateRoute updates an existing route
func (h *RouteController) UpdateRoute(ctx *gin.Context) {
	var routeUpdate dto.RouteCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&routeUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Route ID is required", nil))
		return
	}

	route := dto.ToRouteEntity(&routeUpdate)

	if err := h.RouteUsecase.UpdateRoute(ctx, uint(id), &route); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update route", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Route updated successfully", nil))
}

// DeleteRoute deletes a route by ID
func (h *RouteController) DeleteRoute(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	if err := h.RouteUsecase.DeleteRoute(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete route", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Route deleted successfully", nil))
}
