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
	c := &QuotaController{
		Log:          log,
		Validate:     validate,
		QuotaUsecase: quota_usecase,
	}

	router.GET("/quotas", c.GetAllQuotas)
	router.GET("/quota/:id", c.GetQuotaByID)

	protected.POST("/quota/create", c.CreateQuota)
	protected.POST("/quota/create/bulk", c.CreateQuotaBulk)
	protected.PUT("/quota/update/:id", c.UpdateQuota)
	protected.DELETE("/quota/:id", c.DeleteQuota)

}

func (c *QuotaController) CreateQuota(ctx *gin.Context) {
	request := new(model.WriteQuotaRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.QuotaUsecase.CreateQuota(ctx, request); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}
		c.Log.WithError(err).Error("failed to create Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create Quota", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Quota created successfully", nil))
}

func (c *QuotaController) CreateQuotaBulk(ctx *gin.Context) {
	requests := []*model.WriteQuotaRequest{}

	if err := ctx.ShouldBindJSON(&requests); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body for bulk create")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	for i, req := range requests {
		if err := c.Validate.Struct(req); err != nil {
			c.Log.WithError(err).Error("failed to validate request body at index %d", i)
			errors := validator.ParseErrors(err)
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", map[string]interface{}{
				"index":  i,
				"errors": errors,
			}))
			return
		}
	}

	if err := c.QuotaUsecase.CreateQuotaBulk(ctx, requests); err != nil {
		c.Log.WithError(err).Error("failed to create Quotas in bulk")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create Quotas in bulk", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Quotas created successfully", nil))
}

func (c *QuotaController) GetAllQuotas(ctx *gin.Context) {
	params := response.GetParams(ctx)

	datas, total, err := c.QuotaUsecase.ListQuotas(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve Quotas")
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

func (c *QuotaController) GetQuotaByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid Quota ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid Quota ID", err.Error()))
		return
	}

	data, err := c.QuotaUsecase.GetQuotaByID(ctx, uint(id))

	if err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve Quota", err.Error()))
		return
	}

	if data == nil {
		c.Log.WithField("id", id).Warn("Quota not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Quota not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Quota retrieved successfully", nil))
}

func (c *QuotaController) UpdateQuota(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid or missing Quota ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing Quota ID", nil))
		return
	}

	request := new(model.UpdateQuotaRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.QuotaUsecase.UpdateQuota(ctx, request); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("Quota not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Quota not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("Quota already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("Quota already exists", nil))
			return
		}
		c.Log.WithError(err).WithField("id", id).Error("failed to update Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update Quota", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Quota updated successfully", nil))
}

func (c *QuotaController) DeleteQuota(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid Quota ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid Quota ID", err.Error()))
		return
	}

	if err := c.QuotaUsecase.DeleteQuota(ctx, uint(id)); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to delete Quota")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete Quota", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Quota deleted successfully", nil))
}
