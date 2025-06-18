package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/harbor"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HarborController struct {
	HarborUsecase *harbor.HarborUsecase
	Authenticate  *middleware.AuthenticateMiddleware
	Authorized    *middleware.AuthorizeMiddleware
}

func NewHarborController(
	g *gin.Engine, harbor_usecase *harbor.HarborUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	hc := &HarborController{
		HarborUsecase: harbor_usecase,
		Authenticate:  authtenticate,
		Authorized:    authorized}

	public := g.Group("") // No middleware
	public.GET("/harbors", hc.GetAllHarbors)
	public.GET("/harbor/:id", hc.GetHarborByID)

	protected := g.Group("")
	protected.Use(hc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/harbor/create", hc.CreateHarbor)
	protected.PUT("/harbor/update/:id", hc.UpdateHarbor)
	protected.DELETE("/harbor/:id", hc.DeleteHarbor)
}

func (hc *HarborController) CreateHarbor(ctx *gin.Context) {
	request := new(model.WriteHarborRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := hc.HarborUsecase.CreateHarbor(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Harbor created successfully", nil))
}

func (hc *HarborController) GetAllHarbors(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := hc.HarborUsecase.GetAllHarbors(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbors", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Harbors retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Harbors retrieved successfully", total, params.Limit, params.Page))
}

func (hc *HarborController) GetHarborByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	data, err := hc.HarborUsecase.GetHarborByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbor", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Harbor not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Harbor retrieved successfully", nil))
}

func (hc *HarborController) UpdateHarbor(ctx *gin.Context) {
	request := new(model.UpdateHarborRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Harbor ID is required", nil))
		return
	}

	request.ID = uint(id)
	if err := hc.HarborUsecase.UpdateHarbor(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor updated successfully", nil))
}

func (hc *HarborController) DeleteHarbor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	if err := hc.HarborUsecase.DeleteHarbor(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor deleted successfully", nil))
}
