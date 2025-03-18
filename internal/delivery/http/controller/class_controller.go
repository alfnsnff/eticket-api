package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassController struct {
	ClassUsecase usecase.ClassUsecase
}

// CreateClass handles creating a new class
func (h *ClassController) CreateClass(c *gin.Context) {
	var class entities.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := h.ClassUsecase.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Class created successfully", nil))
}

// GetAllClasses handles retrieving all classes
func (h *ClassController) GetAllClasses(c *gin.Context) {
	classes, err := h.ClassUsecase.GetAllClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve classes", err.Error()))
		return
	}

	classDTOs := dto.ToClassDTOs(classes)
	c.JSON(http.StatusOK, response.NewSuccessResponse(classDTOs, "Classes retrieved successfully", nil))
}

// GetClassByID handles retrieving a class by its ID
func (h *ClassController) GetClassByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	class, err := h.ClassUsecase.GetClassByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve class", err.Error()))
		return
	}

	if class == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Class not found", nil))
		return
	}

	classDTO := dto.ToClassDTO(class)
	c.JSON(http.StatusOK, response.NewSuccessResponse(classDTO, "Class retrieved successfully", nil))
}

// UpdateClass handles updating an existing class
func (h *ClassController) UpdateClass(c *gin.Context) {
	var class entities.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if class.ID == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Class ID is required", nil))
		return
	}

	if err := h.ClassUsecase.UpdateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update class", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class updated successfully", nil))
}

// DeleteClass handles deleting a class by its ID
func (h *ClassController) DeleteClass(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid class ID", err.Error()))
		return
	}

	if err := h.ClassUsecase.DeleteClass(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete class", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Class deleted successfully", nil))
}

// NewClassController creates a new instance of the ClassController.
func NewClassController(classUsecase usecase.ClassUsecase) *ClassController {
	return &ClassController{ClassUsecase: classUsecase}
}
