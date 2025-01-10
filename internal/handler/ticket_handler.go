package handler

import (
    "net/http"
    "eticket-api/internal/domain"
    "eticket-api/internal/service"

    "github.com/gin-gonic/gin"
)

type TicketHandler struct {
    Service *service.TicketService
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
    var ticket domain.Ticket
    if err := c.ShouldBindJSON(&ticket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Service.CreateTicket(&ticket); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "ticket created"})
}

func (h *TicketHandler) GetAllTickets(c *gin.Context) {
    tickets, err := h.Service.GetAllTickets()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, tickets)
}

// Additional handlers for GetByID, Update, Delete
