package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuotaController struct {
	Log          logger.Logger
	Validate     validator.Validator
	QuotaUsecase *usecase.QuotaUsecase
}

func NewQuotaController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	quota_usecase *usecase.QuotaUsecase,

) {
	ac := &QuotaController{
		Log:          log,
		Validate:     validate,
		QuotaUsecase: quota_usecase,
	}

	router.GET("/quotas", ac.GetAllQuotas)
	router.GET("/quota/:id", ac.GetQuotaByID)

	protected.POST("/quota/create", ac.CreateQuota)
	protected.PUT("/quota/update/:id", ac.UpdateQuota)
	protected.DELETE("/quota/:id", ac.DeleteQuota)

}
func (mc *QuotaController) CreateQuota(ctx *gin.Context) {
	request := new(model.WriteQuotaRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		mc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := mc.Validate.Struct(request); err != nil {
		mc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := mc.QuotaUsecase.CreateQuota(ctx, request); err != nil {
		mc.Log.WithError(err).Error("failed to create Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create Quota", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Quota created successfully", nil))
}

func (mc *QuotaController) GetAllQuotas(ctx *gin.Context) {
	params := response.GetParams(ctx)

	datas, total, err := mc.QuotaUsecase.ListQuotas(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		mc.Log.WithError(err).Error("failed to retrieve Quotas")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve Quotas", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Quotas retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (mc *QuotaController) GetQuotaByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		mc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid Quota ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid Quota ID", err.Error()))
		return
	}

	data, err := mc.QuotaUsecase.GetQuotaByID(ctx, uint(id))

	if err != nil {
		mc.Log.WithError(err).WithField("id", id).Error("failed to retrieve Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve Quota", err.Error()))
		return
	}

	if data == nil {
		mc.Log.WithField("id", id).Warn("Quota not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Quota not found", nil))
		return
	}

	mc.Log.WithField("id", id).Info("Quota retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Quota retrieved successfully", nil))
}

func (mc *QuotaController) UpdateQuota(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		mc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid or missing Quota ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing Quota ID", nil))
		return
	}

	request := new(model.UpdateQuotaRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		mc.Log.WithError(err).Error("failed to bind JSON request body")
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

	if err := mc.QuotaUsecase.UpdateQuota(ctx, request); err != nil {
		mc.Log.WithError(err).WithField("id", id).Error("failed to update Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update Quota", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Quota updated successfully", nil))
}

func (mc *QuotaController) DeleteQuota(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		mc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid Quota ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid Quota ID", err.Error()))
		return
	}

	if err := mc.QuotaUsecase.DeleteQuota(ctx, uint(id)); err != nil {
		mc.Log.WithError(err).WithField("id", id).Error("failed to delete Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete Quota", err.Error()))
		return
	}

	mc.Log.WithField("id", id).Info("Quota deleted successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Quota deleted successfully", nil))
}
