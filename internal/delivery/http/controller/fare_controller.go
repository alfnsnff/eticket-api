package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/model" // Import the response package
	"eticket-api/internal/usecase/fare"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FareController struct {
	FareUsecase  *fare.FareUsecase
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewFareController(
	g *gin.Engine, Fare_usecase *fare.FareUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	fc := &FareController{
		FareUsecase:  Fare_usecase,
		Authenticate: authtenticate,
		Authorized:   authorized}

	public := g.Group("") // No middleware
	public.GET("/fares", fc.GetAllFares)
	public.GET("/fare/:id", fc.GetFareByID)

	protected := g.Group("")
	protected.Use(fc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/fare/create", fc.CreateFare)
	protected.PUT("/fare/update/:id", fc.UpdateFare)
	protected.DELETE("/fare/:id", fc.DeleteFare)
}

func (fc *FareController) CreateFare(ctx *gin.Context) {
	request := new(model.WriteFareRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
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
	request := new(model.UpdateFareRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Fare ID is required", nil))
		return
	}
	request.ID = uint(id)

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
