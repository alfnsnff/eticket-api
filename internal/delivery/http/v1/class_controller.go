package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"

	"eticket-api/internal/delivery/http/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassController struct {
	Validate     validator.Validator
	Log          logger.Logger
	ClassUsecase *usecase.ClassUsecase
}

func NewClassController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	class_usecase *usecase.ClassUsecase,

) {
	cc := &ClassController{
		Log:          log,
		Validate:     validate,
		ClassUsecase: class_usecase,
	}

	router.GET("/classes", cc.GetAllClasses)
	router.GET("/class/:id", cc.GetClassByID)

	protected.POST("/class/create", cc.CreateClass)
	protected.PUT("/class/update/:id", cc.UpdateClass)
	protected.DELETE("/class/:id", cc.DeleteClass)
}

func (cc *ClassController) CreateClass(ctx *gin.Context) {
	request := new(model.WriteClassRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		cc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := cc.Validate.Struct(request); err != nil {
		cc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := cc.ClassUsecase.CreateClass(ctx, request); err != nil {
		cc.Log.WithError(err).Error("failed to create class")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Class created successfully", nil))
}

func (cc *ClassController) GetAllClasses(ctx *gin.Context) {

	params := response.GetParams(ctx)
	datas, total, err := cc.ClassUsecase.ListClasses(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		cc.Log.WithError(err).Error("failed to retrieve classes")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve classes", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Classes retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (cc *ClassController) GetClassByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		cc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse class ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	data, err := cc.ClassUsecase.GetClassByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			cc.Log.WithField("id", id).Warn("class not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
			return
		}

		cc.Log.WithError(err).WithField("id", id).Error("failed to retrieve class")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve class", err.Error()))
		return
	}

	if data == nil {
		cc.Log.WithField("id", id).Warn("class not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Class retrieved successfully", nil))
}

func (cc *ClassController) UpdateClass(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		cc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse class ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing class ID", nil))
		return
	}

	request := new(model.UpdateClassRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		cc.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := cc.Validate.Struct(request); err != nil {
		cc.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := cc.ClassUsecase.UpdateClass(ctx, request); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			cc.Log.WithField("id", id).Warn("class not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
			return
		}

		cc.Log.WithError(err).WithField("id", id).Error("failed to update class")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update class", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class updated successfully", nil))
}

func (cc *ClassController) DeleteClass(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		cc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse class ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	if err := cc.ClassUsecase.DeleteClass(ctx, uint(id)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			cc.Log.WithField("id", id).Warn("class not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
			return
		}

		cc.Log.WithError(err).WithField("id", id).Error("failed to delete class")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete class", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class deleted successfully", nil))
}
