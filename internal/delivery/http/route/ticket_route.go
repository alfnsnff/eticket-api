package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type TicketRouter struct {
	Controller   *controller.TicketController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewTicketRouter(ticket_controller *controller.TicketController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *TicketRouter {
	return &TicketRouter{Controller: ticket_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i TicketRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	tc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/tickets", tc.GetAllTickets)
	public.GET("/ticket/:id", tc.GetTicketByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/ticket/create", tc.CreateTicket)
	protected.PUT("/ticket//update:id", tc.UpdateTicket)
	protected.DELETE("/ticket/:id", tc.DeleteTicket)
}
