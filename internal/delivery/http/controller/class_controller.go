package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassController struct {
	ClassUsecase *usecase.ClassUsecase
}

// NewClassController creates a new instance of the ClassController.
func NewClassController(classUsecase *usecase.ClassUsecase) *ClassController {
	return &ClassController{ClassUsecase: classUsecase}
}

// CreateClass handles creating a new class
func (h *ClassController) CreateClass(ctx *gin.Context) {
	var classCreate dto.ClassCreate

	if err := ctx.ShouldBindJSON(&classCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	class := dto.ToClassEntity(&classCreate)

	if err := h.ClassUsecase.CreateClass(ctx, &class); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Class created successfully", nil))
}

// GetAllClasses handles retrieving all classes
func (h *ClassController) GetAllClasses(ctx *gin.Context) {
	classes, err := h.ClassUsecase.GetAllClasses(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve classes", err.Error()))
		return
	}

	classDTOs := dto.ToClassDTOs(classes)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(classDTOs, "Classes retrieved successfully", nil))
}

// GetClassByID handles retrieving a class by its ID
func (h *ClassController) GetClassByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	class, err := h.ClassUsecase.GetClassByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve class", err.Error()))
		return
	}

	if class == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
		return
	}

	classDTO := dto.ToClassDTO(class)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(classDTO, "Class retrieved successfully", nil))
}

// UpdateClass handles updating an existing class
func (h *ClassController) UpdateClass(ctx *gin.Context) {
	var classUpdate dto.ClassCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&classUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Class ID is required", nil))
		return
	}

	class := dto.ToClassEntity(&classUpdate)

	if err := h.ClassUsecase.UpdateClass(ctx, uint(id), &class); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update class", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class updated successfully", nil))
}

// DeleteClass handles deleting a class by its ID
func (h *ClassController) DeleteClass(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	if err := h.ClassUsecase.DeleteClass(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete class", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class deleted successfully", nil))
}
