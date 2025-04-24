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
	TicketUsecase usecase.TicketUsecase
}

func (h *TicketController) ValidateTicket(c *gin.Context) {
	var req dto.TicketSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	res, err := h.TicketUsecase.ValidateTicketSelection(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(res, "Ticket availability validated", nil))
}

func (h *TicketController) GetBookedCount(c *gin.Context) {
	var req dto.TicketBookedCount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}
	res, err := h.TicketUsecase.GetBookedCount(req.ScheduleID, req.PriceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(res, "Ticket availability counted", nil))
}

// CreateTicket handles creating a new ticket
func (h *TicketController) CreateTicket(c *gin.Context) {
	var ticketCreate dto.TicketCreate
	if err := c.ShouldBindJSON(&ticketCreate); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	ticket := dto.ToTicketEntity(&ticketCreate)

	if err := h.TicketUsecase.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ticket", err.Error())) // Use response.
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ticket created successfully", nil)) // Use response.
}

// GetAllTickets handles retrieving all tickets
func (h *TicketController) GetAllTickets(c *gin.Context) {
	tickets, err := h.TicketUsecase.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error())) // Use response.
		return
	}

	ticketDTOs := dto.ToTicketDTOs(tickets)
	c.JSON(http.StatusOK, response.NewSuccessResponse(ticketDTOs, "Tickets retrieved successfully", nil)) // Use response.
}

// GetTicketByID handles retrieving a ticket by its ID
func (h *TicketController) GetTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	ticket, err := h.TicketUsecase.GetTicketByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ticket", err.Error())) // Use response.
		return
	}

	if ticket == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Ticket not found", nil)) // Use response.
		return
	}

	ticketDTO := dto.ToTicketDTO(ticket)
	c.JSON(http.StatusOK, response.NewSuccessResponse(ticketDTO, "Ticket retrieved successfully", nil))
}

// UpdateTicket handles updating an existing ticket
func (h *TicketController) UpdateTicket(c *gin.Context) {
	var ticketUpdate dto.TicketCreate

	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&ticketUpdate); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error())) // Use response.
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Ticket ID is required", nil)) // Use response.
		return
	}

	ticket := dto.ToTicketEntity(&ticketUpdate)

	if err := h.TicketUsecase.UpdateTicket(uint(id), &ticket); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ticket", err.Error())) // Use response.
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket updated successfully", nil)) // Use response.
}

// DeleteTicket handles deleting a ticket by its ID
func (h *TicketController) DeleteTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error())) // Use response.
		return
	}

	if err := h.TicketUsecase.DeleteTicket(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ticket", err.Error())) // Use response.
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket deleted successfully", nil)) // Use response.
}

// NewTicketController creates a new instance of TicketController.
func NewTicketController(ticketUsecase usecase.TicketUsecase) *TicketController {
	return &TicketController{TicketUsecase: ticketUsecase}
}
