package controller

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	TicketUsecase usecase.TicketUsecase
}

// CreateTicket handles creating a new ticket
func (h *TicketController) CreateTicket(c *gin.Context) {
	var ticket domain.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TicketUsecase.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ticket created"})
}

// GetAllTickets handles retrieving all tickets
func (h *TicketController) GetAllTickets(c *gin.Context) {
	tickets, err := h.TicketUsecase.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// GetTicketByID handles retrieving a ticket by its ID
func (h *TicketController) GetTicketByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	ticket, err := h.TicketUsecase.GetTicketByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// UpdateTicket handles updating an existing ticket
func (h *TicketController) UpdateTicket(c *gin.Context) {
	var ticket domain.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ticket.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket ID is required"})
		return
	}

	if err := h.TicketUsecase.UpdateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket updated"})
}

// DeleteTicket handles deleting a ticket by its ID
func (h *TicketController) DeleteTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.TicketUsecase.DeleteTicket(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket deleted"})
}
