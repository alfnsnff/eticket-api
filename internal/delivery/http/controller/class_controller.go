package controller

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassController struct {
	ClassUsecase usecase.ClassUsecase
}

// CreateClass handles creating a new class
func (h *ClassController) CreateClass(c *gin.Context) {
	var class domain.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ClassUsecase.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "class created"})
}

// GetAllClasses handles retrieving all classes
func (h *ClassController) GetAllClasses(c *gin.Context) {
	classes, err := h.ClassUsecase.GetAllClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, classes)
}

// GetClassByID handles retrieving a class by its ID
func (h *ClassController) GetClassByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	class, err := h.ClassUsecase.GetClassByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if class == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}

	c.JSON(http.StatusOK, class)
}

// UpdateClass handles updating an existing class
func (h *ClassController) UpdateClass(c *gin.Context) {
	var class domain.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if class.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class ID is required"})
		return
	}

	if err := h.ClassUsecase.UpdateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "class updated"})
}

// DeleteClass handles deleting a class by its ID
func (h *ClassController) DeleteClass(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.ClassUsecase.DeleteClass(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "class deleted"})
}
