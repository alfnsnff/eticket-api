package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/manifest" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManifestController struct {
	Log             logger.Logger
	Validate        validator.Validator
	ManifestUsecase *manifest.ManifestUsecase
}

func NewManifestController(
	log logger.Logger,
	validate validator.Validator,
	manifest_usecase *manifest.ManifestUsecase,
) *ManifestController {
	return &ManifestController{
		Log:             log,
		Validate:        validate,
		ManifestUsecase: manifest_usecase,
	}
}

func (mc *ManifestController) CreateManifest(ctx *gin.Context) {
	request := new(model.WriteManifestRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := mc.Validate.Struct(request); err != nil {
		mc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateManifestRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := mc.Validate.Struct(request); err != nil {
		mc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

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
