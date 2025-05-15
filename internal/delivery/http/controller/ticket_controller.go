package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/helper/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	TicketUsecase *usecase.TicketUsecase
}

func NewTicketController(ticket_usecase *usecase.TicketUsecase) *TicketController {
	return &TicketController{TicketUsecase: ticket_usecase}
}

func (tc *TicketController) CreateTicket(ctx *gin.Context) {
	request := new(model.WriteTicketRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	if err := tc.TicketUsecase.CreateTicket(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ticket created successfully", nil)) // Use response.
}

func (tc *TicketController) GetAllTickets(ctx *gin.Context) {
	datas, err := tc.TicketUsecase.GetAllTickets(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Tickets retrieved successfully", nil)) // Use response.
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
	request := new(model.UpdateTicketRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ticket ID is required", nil)) // Use response.
		return
	}

	if err := tc.TicketUsecase.UpdateTicket(ctx, uint(id), request); err != nil {
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
