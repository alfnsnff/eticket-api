package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManifestController struct {
	ManifestUsecase *usecase.ManifestUsecase
}

func NewManifestController(manifest_usecase *usecase.ManifestUsecase) *ManifestController {
	return &ManifestController{ManifestUsecase: manifest_usecase}
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
	datas, err := mc.ManifestUsecase.GetAllManifests(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve manifests", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Manifests retrieved successfully", nil))
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

	if err := mc.ManifestUsecase.UpdateManifest(ctx, uint(id), request); err != nil {
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
