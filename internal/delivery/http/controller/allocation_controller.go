package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/allocation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AllocationController struct {
	Log               logger.Logger
	Validate          validator.Validator
	AllocationUsecase *allocation.AllocationUsecase
	Authenticate      *middleware.AuthenticateMiddleware
	Authorized        *middleware.AuthorizeMiddleware
}

func NewAllocationController(
	router *gin.Engine,
	log logger.Logger,
	validate validator.Validator,
	allocation_usecase *allocation.AllocationUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	ac := &AllocationController{AllocationUsecase: allocation_usecase,
		Log:          log,
		Validate:     validate,
		Authenticate: authtenticate,
		Authorized:   authorized,
	}

	public := router.Group("/api/v1") // No middleware
	public.GET("/allocations", ac.GetAllAllocations)
	public.GET("/allocation/:id", ac.GetAllocationByID)

	protected := router.Group("/api/v1")
	protected.Use(ac.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	protected.POST("/allocation/create", ac.CreateAllocation)
	protected.PUT("/allocation/update/:id", ac.UpdateAllocation)
	protected.DELETE("/allocation/:id", ac.DeleteAllocation)

}
func (mc *AllocationController) CreateAllocation(ctx *gin.Context) {
	request := new(model.WriteAllocationRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := mc.Validate.Struct(request); err != nil {
		mc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
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

	datas, total, err := mc.AllocationUsecase.GetAllAllocations(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve allocations", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Allocations retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Allocations retrieved successfully", total, params.Limit, params.Page))
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateAllocationRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
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

	if err := mc.AllocationUsecase.UpdateAllocation(ctx, request); err != nil {
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
