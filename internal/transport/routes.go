package transport

import (
    "eticket-api/internal/handler"

    "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ticketHandler *handler.TicketHandler) {
    v1 := router.Group("/api/v1")
    {
        tickets := v1.Group("/tickets")
        {
            tickets.POST("/", ticketHandler.CreateTicket)
            tickets.GET("/", ticketHandler.GetAllTickets)
            // Other routes (GET /:id, PUT /:id, DELETE /:id)
        }
    }
}
