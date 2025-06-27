package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HarborController struct {
	Validate      validator.Validator
	Log           logger.Logger
	HarborUsecase *usecase.HarborUsecase
}

func NewHarborController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	harbor_usecase *usecase.HarborUsecase,

) {
	hc := &HarborController{
		Log:           log,
		Validate:      validate,
		HarborUsecase: harbor_usecase,
	}
	router.GET("/harbors", hc.GetAllHarbors)
	router.GET("/harbor/:id", hc.GetHarborByID)

	protected.POST("/harbor/create", hc.CreateHarbor)
	protected.PUT("/harbor/update/:id", hc.UpdateHarbor)
	protected.DELETE("/harbor/:id", hc.DeleteHarbor)
}

func (hc *HarborController) CreateHarbor(ctx *gin.Context) {
	request := new(model.WriteHarborRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		hc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := hc.Validate.Struct(request); err != nil {
		hc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := hc.HarborUsecase.CreateHarbor(ctx, request); err != nil {
		hc.Log.WithError(err).Error("failed to create harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Harbor created successfully", nil))
}

func (hc *HarborController) GetAllHarbors(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := hc.HarborUsecase.ListHarbors(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		hc.Log.WithError(err).Error("failed to retrieve harbors")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbors", err.Error()))
		return
	}

	hc.Log.WithField("count", total).Info("Harbors retrieved successfully")
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
}

func (hc *HarborController) GetHarborByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		hc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse harbor ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	data, err := hc.HarborUsecase.GetHarborByID(ctx, uint(id))

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			hc.Log.WithField("id", id).Warn("harbor not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("class not found", nil))
			return
		}

		hc.Log.WithError(err).WithField("id", id).Error("failed to retrieve harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbor", err.Error()))
		return
	}

	if data == nil {
		hc.Log.WithField("id", id).Warn("harbor not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Harbor not found", nil))
		return
	}
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Harbor retrieved successfully", nil))
}

func (hc *HarborController) UpdateHarbor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		hc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse harbor ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing harbor ID", nil))
		return
	}

	request := new(model.UpdateHarborRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		hc.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := hc.Validate.Struct(request); err != nil {
		hc.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := hc.HarborUsecase.UpdateHarbor(ctx, request); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			hc.Log.WithField("id", id).Warn("harbor not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("harbor not found", nil))
			return
		}

		hc.Log.WithError(err).WithField("id", id).Error("failed to update harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor updated successfully", nil))
}

func (hc *HarborController) DeleteHarbor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		hc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse harbor ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	if err := hc.HarborUsecase.DeleteHarbor(ctx, uint(id)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			hc.Log.WithField("id", id).Warn("harbor not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("harbor not found", nil))
			return
		}

		hc.Log.WithError(err).WithField("id", id).Error("failed to delete harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor deleted successfully", nil))
}
