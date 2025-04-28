package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	TicketUsecase *usecase.TicketUsecase
}

// NewTicketController creates a new instance of TicketController.
func NewTicketController(ticketUsecase *usecase.TicketUsecase) *TicketController {
	return &TicketController{TicketUsecase: ticketUsecase}
}

func (h *TicketController) ValidateTicket(ctx *gin.Context) {
	var req dto.TicketSelectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	res, err := h.TicketUsecase.ValidateTicketSelection(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(res, "Ticket availability validated", nil))
}

func (h *TicketController) GetBookedCount(ctx *gin.Context) {
	var req dto.TicketBookedCount
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}
	res, err := h.TicketUsecase.GetBookedCount(ctx, req.ScheduleID, req.PriceID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(res, "Ticket availability counted", nil))
}

// CreateTicket handles creating a new ticket
func (h *TicketController) CreateTicket(ctx *gin.Context) {
	var ticketCreate dto.TicketCreate
	if err := ctx.ShouldBindJSON(&ticketCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	ticket := dto.ToTicketEntity(&ticketCreate)

	if err := h.TicketUsecase.CreateTicket(ctx, &ticket); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ticket created successfully", nil)) // Use response.
}

// GetAllTickets handles retrieving all tickets
func (h *TicketController) GetAllTickets(ctx *gin.Context) {
	tickets, err := h.TicketUsecase.GetAllTickets(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error())) // Use response.
		return
	}

	ticketDTOs := dto.ToTicketDTOs(tickets)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(ticketDTOs, "Tickets retrieved successfully", nil)) // Use response.
}

// GetTicketByID handles retrieving a ticket by its ID
func (h *TicketController) GetTicketByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	ticket, err := h.TicketUsecase.GetTicketByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ticket", err.Error())) // Use response.
		return
	}

	if ticket == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ticket not found", nil)) // Use response.
		return
	}

	ticketDTO := dto.ToTicketDTO(ticket)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(ticketDTO, "Ticket retrieved successfully", nil))
}

// UpdateTicket handles updating an existing ticket
func (h *TicketController) UpdateTicket(ctx *gin.Context) {
	var ticketUpdate dto.TicketCreate

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&ticketUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Ticket ID is required", nil)) // Use response.
		return
	}

	ticket := dto.ToTicketEntity(&ticketUpdate)

	if err := h.TicketUsecase.UpdateTicket(ctx, uint(id), &ticket); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket updated successfully", nil)) // Use response.
}

// DeleteTicket handles deleting a ticket by its ID
func (h *TicketController) DeleteTicket(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	if err := h.TicketUsecase.DeleteTicket(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ticket", err.Error())) // Use response.
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket deleted successfully", nil)) // Use response.
}
