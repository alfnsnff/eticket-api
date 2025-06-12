package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/allocation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AllocationController struct {
	AllocationUsecase *allocation.AllocationUsecase
}

func NewAllocationController(allocation_usecase *allocation.AllocationUsecase) *AllocationController {
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
	params := response.GetParams(ctx)

	datas, total, err := mc.AllocationUsecase.GetAllAllocations(ctx, params.Limit, params.Offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve allocations", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Allocations retrieved successfully", total, params.Limit, params.Page))
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
