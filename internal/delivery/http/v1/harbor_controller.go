package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	requests "eticket-api/internal/delivery/http/v1/request"
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
	c := &HarborController{
		Log:           log,
		Validate:      validate,
		HarborUsecase: harbor_usecase,
	}
	router.GET("/harbors", c.GetAllHarbors)
	router.GET("/harbor/:id", c.GetHarborByID)

	protected.POST("/harbor/create", c.CreateHarbor)
	protected.PUT("/harbor/update/:id", c.UpdateHarbor)
	protected.DELETE("/harbor/:id", c.DeleteHarbor)
}

func (c *HarborController) CreateHarbor(ctx *gin.Context) {
	request := new(requests.CreateHarborRequest)

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

	if err := c.HarborUsecase.CreateHarbor(ctx, requests.HarborFromCreate(request)); err != nil {

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}
		c.Log.WithError(err).Error("failed to create harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Harbor created successfully", nil))
}

func (c *HarborController) GetAllHarbors(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := c.HarborUsecase.ListHarbors(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve harbors")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbors", err.Error()))
		return
	}

	responses := make([]*requests.HarborResponse, len(datas))
	for i, data := range datas {
		responses[i] = requests.HarborToResponse(data)
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		responses,
		"Harbors retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (c *HarborController) GetHarborByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse harbor ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	data, err := c.HarborUsecase.GetHarborByID(ctx, uint(id))

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("harbor not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("class not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(requests.HarborToResponse(data), "Harbor retrieved successfully", nil))
}

func (c *HarborController) UpdateHarbor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse harbor ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing harbor ID", nil))
		return
	}

	request := new(requests.UpdateHarborRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.HarborUsecase.UpdateHarbor(ctx, requests.HarborFromUpdate(request)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("harbor not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("harbor not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("harbor already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("harbor already exists", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to update harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor updated successfully", nil))
}

func (c *HarborController) DeleteHarbor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse harbor ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid harbor ID", err.Error()))
		return
	}

	if err := c.HarborUsecase.DeleteHarbor(ctx, uint(id)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("harbor not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("harbor not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to delete harbor")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete harbor", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Harbor deleted successfully", nil))
}
