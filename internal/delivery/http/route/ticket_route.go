package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewTicketRouter(ic *injector.Container, rg *gin.RouterGroup) {
	tr := repository.NewTicketRepository()
	scr := repository.NewScheduleRepository()
	fr := repository.NewFareRepository()
	sr := repository.NewSessionRepository()

	tc := &controller.TicketController{
		TicketUsecase: usecase.NewTicketUsecase(ic.Tx, tr, scr, fr, sr),
	}

	public := rg.Group("") // No middleware
	public.GET("/tickets", tc.GetAllTickets)
	public.GET("/ticket/:id", tc.GetTicketByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.Repository.UserRepository, ic.Repository.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/ticket/create", tc.CreateTicket)
	protected.PUT("/ticket//update:id", tc.UpdateTicket)
	protected.DELETE("/ticket/:id", tc.DeleteTicket)
}
