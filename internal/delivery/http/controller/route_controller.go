package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RouteController struct {
	RouteUsecase usecase.RouteUsecase
}

// CreateRoute creates a new route
func (h *RouteController) CreateRoute(c *gin.Context) {
	var route entities.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := h.RouteUsecase.CreateRoute(&route); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create route", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Route created successfully", nil))
}

// GetAllRoutes retrieves all routes
func (h *RouteController) GetAllRoutes(c *gin.Context) {
	routes, err := h.RouteUsecase.GetAllRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve routes", err.Error()))
		return
	}

	routeDTOs := dto.ToRouteDTOs(routes)
	c.JSON(http.StatusOK, response.NewSuccessResponse(routeDTOs, "Routes retrieved successfully", nil))
}

// GetRouteByID retrieves a single route by ID
func (h *RouteController) GetRouteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	route, err := h.RouteUsecase.GetRouteByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve route", err.Error())) // More specific error
		return
	}
	if route == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Route not found", nil))
		return
	}

	routeDTO := dto.ToRouteDTO(route)
	c.JSON(http.StatusOK, response.NewSuccessResponse(routeDTO, "Route retrieved successfully", nil))
}

// UpdateRoute updates an existing route
func (h *RouteController) UpdateRoute(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	var route entities.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	route.ID = uint(id) // Ensure the ID is set for updating
	if err := h.RouteUsecase.UpdateRoute(&route); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update route", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Route updated successfully", nil))
}

// DeleteRoute deletes a route by ID
func (h *RouteController) DeleteRoute(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid route ID", err.Error()))
		return
	}

	if err := h.RouteUsecase.DeleteRoute(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete route", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Route deleted successfully", nil))
}

// NewRouteController creates a new RouteController instance.
func NewRouteController(routeUsecase usecase.RouteUsecase) *RouteController {
	return &RouteController{RouteUsecase: routeUsecase}
}
