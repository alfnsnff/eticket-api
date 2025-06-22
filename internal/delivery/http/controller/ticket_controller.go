package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/ticket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	Validate      validator.Validator
	Log           logger.Logger
	TicketUsecase *ticket.TicketUsecase
	Authenticate  *middleware.AuthenticateMiddleware
	Authorized    *middleware.AuthorizeMiddleware
}

func NewTicketController(
	g *gin.Engine,
	log logger.Logger,
	validate validator.Validator,
	ticket_usecase *ticket.TicketUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	tc := &TicketController{
		TicketUsecase: ticket_usecase,
		Authenticate:  authtenticate,
		Authorized:    authorized,
		Validate:      validate,
		Log:           log,
	}

	public := g.Group("/api/v1") // No middleware
	public.GET("/tickets", tc.GetAllTickets)
	public.GET("/ticket/:id", tc.GetTicketByID)

	protected := g.Group("/api/v1")
	protected.Use(tc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())

	public.GET("/ticket/schedule/:id", tc.GetAllTicketsByScheduleID)
	protected.POST("/ticket/check-in/:id", tc.CheckIn)
	protected.POST("/ticket/create", tc.CreateTicket)
	protected.PUT("/ticket//update:id", tc.UpdateTicket)
	protected.DELETE("/ticket/:id", tc.DeleteTicket)
}

func (tc *TicketController) CreateTicket(ctx *gin.Context) {
	request := new(model.WriteTicketRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	if err := tc.Validate.Struct(request); err != nil {
		tc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := tc.TicketUsecase.CreateTicket(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ticket created successfully", nil)) // Use response.
}

func (tc *TicketController) GetAllTickets(ctx *gin.Context) {
	params := response.GetParams(ctx)

	datas, total, err := tc.TicketUsecase.GetAllTickets(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Tickets retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Tickets retrieved successfully", total, params.Limit, params.Page))
}

func (tc *TicketController) GetAllTicketsByScheduleID(ctx *gin.Context) {
	params := response.GetParams(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	datas, total, err := tc.TicketUsecase.GetAllTicketsByScheduleID(ctx, id, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Tickets retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Tickets retrieved successfully", total, params.Limit, params.Page))
}

func (tc *TicketController) GetTicketByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	data, err := tc.TicketUsecase.GetTicketByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ticket", err.Error())) // Use response.
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ticket not found", nil)) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ticket retrieved successfully", nil))
}

func (tc *TicketController) UpdateTicket(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ship ID", nil))
		return
	}

	request := new(model.UpdateTicketRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	request.ID = uint(id)
	if err := tc.Validate.Struct(request); err != nil {
		tc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := tc.TicketUsecase.UpdateTicket(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket updated successfully", nil)) // Use response.
}

func (tc *TicketController) DeleteTicket(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	if err := tc.TicketUsecase.DeleteTicket(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket deleted successfully", nil)) // Use response.
}

func (tc *TicketController) CheckIn(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	if err := tc.TicketUsecase.CheckIn(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket deleted successfully", nil)) // Use response.
}
