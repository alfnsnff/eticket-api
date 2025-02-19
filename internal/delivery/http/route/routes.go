package route

import (
	"eticket-api/internal/delivery/http/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ticketController *controller.TicketController, classController *controller.ClassController, routeController *controller.RouteController) {
	v1 := router.Group("/api/v1")
	{
		tickets := v1.Group("/tickets")
		{
			tickets.POST("", ticketController.CreateTicket)
			tickets.GET("", ticketController.GetAllTickets)
			tickets.GET("/:id", ticketController.GetTicketByID)
			tickets.PUT("/:id", ticketController.UpdateTicket)
			tickets.DELETE("/:id", ticketController.DeleteTicket)
		}
		classes := v1.Group("/classes")
		{
			classes.POST("", classController.CreateClass)
			classes.GET("", classController.GetAllClasses)
			classes.GET("/:id", classController.GetClassByID)
			classes.PUT("/:id", classController.UpdateClass)
			classes.DELETE("/:id", classController.DeleteClass)
		}
		routes := v1.Group("/routes")
		{
			routes.POST("", routeController.CreateRoute)
			routes.GET("", routeController.GetAllRoutes)
			routes.GET("/:id", routeController.GetRouteByID)
			routes.PUT("/:id", routeController.UpdateRoute)
			routes.DELETE("/:id", routeController.DeleteRoute)

		}
	}
}
