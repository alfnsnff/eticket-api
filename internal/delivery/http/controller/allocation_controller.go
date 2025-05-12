package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase" // Import the response package
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AllocationController struct {
	AllocationUsecase *usecase.AllocationUsecase
}

func NewAllocationController(allocation_usecase *usecase.AllocationUsecase) *AllocationController {
	return &AllocationController{AllocationUsecase: allocation_usecase}
}

func (mc *AllocationController) CreateAllocation(ctx *gin.Context) {
	request := new(model.WriteAllocationRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := mc.AllocationUsecase.CreateAllocation(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ship", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Allocation created successfully", nil))
}

func (mc *AllocationController) GetAllAllocations(ctx *gin.Context) {
	datas, err := mc.AllocationUsecase.GetAllAllocations(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ships", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Allocations retrieved successfully", nil))
}

func (mc *AllocationController) GetAllocationByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid allocation ID", err.Error()))
		return
	}

	data, err := mc.AllocationUsecase.GetAllocationByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve allocation", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Allocation not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Allocation retrieved successfully", nil))
}

func (mc *AllocationController) UpdateAllocation(ctx *gin.Context) {
	request := new(model.UpdateAllocationRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Allocation ID is required", nil))
		return
	}

	if err := mc.AllocationUsecase.UpdateAllocation(ctx, uint(id), request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update allocation", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Allocation updated successfully", nil))
}

func (mc *AllocationController) DeleteAllocation(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid alocation ID", err.Error()))
		return
	}

	if err := mc.AllocationUsecase.DeleteAllocation(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete allocation", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "allocation deleted successfully", nil))
}
