package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/ship" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipController struct {
	ShipUsecase  *ship.ShipUsecase
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewShipController(
	g *gin.Engine, ship_usecase *ship.ShipUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	shc := &ShipController{
		ShipUsecase:  ship_usecase,
		Authenticate: authtenticate,
		Authorized:   authorized}

	public := g.Group("") // No middleware
	public.GET("/ships", shc.GetAllShips)
	public.GET("/ships/:id", shc.GetShipByID)

	protected := g.Group("")
	protected.Use(shc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/ship/create", shc.CreateShip)
	protected.PUT("/ship/:id", shc.UpdateShip)
	protected.DELETE("/ship/:id", shc.DeleteShip)
}

func (shc *ShipController) CreateShip(ctx *gin.Context) {
	request := new(model.WriteShipRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := shc.ShipUsecase.CreateShip(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ship created successfully", nil))
}

func (shc *ShipController) GetAllShips(ctx *gin.Context) {
	params := response.GetParams(ctx)

	datas, total, err := shc.ShipUsecase.GetAllShips(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

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
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	data, err := shc.ShipUsecase.GetShipByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ship", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ship not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ship retrieved successfully", nil))
}

func (shc *ShipController) UpdateShip(ctx *gin.Context) {
	request := new(model.UpdateShipRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ship ID is required", nil))
		return
	}
	request.ID = uint(id)
	if err := shc.ShipUsecase.UpdateShip(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship updated successfully", nil))
}

func (shc *ShipController) DeleteShip(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ship ID", err.Error()))
		return
	}

	if err := shc.ShipUsecase.DeleteShip(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ship deleted successfully", nil))
}
