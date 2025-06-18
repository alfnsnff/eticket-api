package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type TicketRoute struct {
	Controller   *controller.TicketController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewTicketRoute(ticket_controller *controller.TicketController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *TicketRoute {
	return &TicketRoute{Controller: ticket_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i TicketRoute) Set(router *gin.Engine) {
	tc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/tickets", tc.GetAllTickets)
	public.GET("/ticket/:id", tc.GetTicketByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	protected.Use(i.Authorized.Handle())

	public.GET("/ticket/schedule/:id", tc.GetAllTicketsByScheduleID)
	protected.POST("/ticket/check-in/:id", tc.CheckIn)
	protected.POST("/ticket/create", tc.CreateTicket)
	protected.PUT("/ticket//update:id", tc.UpdateTicket)
	protected.DELETE("/ticket/:id", tc.DeleteTicket)
}
