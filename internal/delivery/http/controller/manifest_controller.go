package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/manifest" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManifestController struct {
	ManifestUsecase *manifest.ManifestUsecase
	Authenticate    *middleware.AuthenticateMiddleware
	Authorized      *middleware.AuthorizeMiddleware
}

func NewManifestController(
	g *gin.Engine, manifest_usecase *manifest.ManifestUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	mc := &ManifestController{
		ManifestUsecase: manifest_usecase,
		Authenticate:    authtenticate,
		Authorized:      authorized}

	public := g.Group("") // No middleware
	public.GET("/manifests", mc.GetAllManifests)
	public.GET("/manifest/:id", mc.GetManifestByID)

	protected := g.Group("")
	protected.Use(mc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/manifest/create", mc.CreateManifest)
	protected.PUT("/manifest/update/:id", mc.UpdateManifest)
	protected.DELETE("/manifest/:id", mc.DeleteManifest)
}

func (mc *ManifestController) CreateManifest(ctx *gin.Context) {
	request := new(model.WriteManifestRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := mc.ManifestUsecase.CreateManifest(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create manifest", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Manifest created successfully", nil))
}

func (mc *ManifestController) GetAllManifests(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := mc.ManifestUsecase.GetAllManifests(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve manifests", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Manifests retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Manifests retrieved successfully", total, params.Limit, params.Page))
}

func (mc *ManifestController) GetManifestByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid manifest ID", err.Error()))
		return
	}

	data, err := mc.ManifestUsecase.GetManifestByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve manifest", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Manifest not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Manifest retrieved successfully", nil))
}

func (mc *ManifestController) UpdateManifest(ctx *gin.Context) {
	request := new(model.UpdateManifestRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Manifest ID is required", nil))
		return
	}
	request.ID = uint(id)
	if err := mc.ManifestUsecase.UpdateManifest(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update manifest", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Manifest updated successfully", nil))
}

func (mc *ManifestController) DeleteManifest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid manifest ID", err.Error()))
		return
	}

	if err := mc.ManifestUsecase.DeleteManifest(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete manifest", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Manifest deleted successfully", nil))
}
