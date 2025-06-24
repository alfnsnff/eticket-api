package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/class"
	"fmt"

	"eticket-api/internal/delivery/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassController struct {
	Log          logger.Logger
	Validate     validator.Validator
	ClassUsecase *class.ClassUsecase
}

func NewClassController(
	log logger.Logger,
	validate validator.Validator,
	class_usecase *class.ClassUsecase,
) *ClassController {
	return &ClassController{
		Log:          log,
		Validate:     validate,
		ClassUsecase: class_usecase,
	}
}

func (cc *ClassController) CreateClass(ctx *gin.Context) {
	request := new(model.WriteClassRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Class created successfully", nil))
}

func (cc *ClassController) GetAllClasses(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := cc.ClassUsecase.GetAllClasses(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve classes", err.Error()))
		return
	}
	fmt.Println("🎯 Hit GetAllTickets")
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

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Classes retrieved successfully", total, params.Limit, params.Page))
}

func (cc *ClassController) GetClassByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	data, err := cc.ClassUsecase.GetClassByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve class", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Class retrieved successfully", nil))
}

func (cc *ClassController) UpdateClass(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateClassRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := cc.Validate.Struct(request); err != nil {
		cc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := cc.ClassUsecase.UpdateClass(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update class", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class updated successfully", nil))
}

func (cc *ClassController) DeleteClass(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	if err := cc.ClassUsecase.DeleteClass(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete class", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class deleted successfully", nil))
}
