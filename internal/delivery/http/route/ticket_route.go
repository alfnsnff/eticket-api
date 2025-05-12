package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTicketRouter(db *gorm.DB, group *gin.RouterGroup) {
	tr := repository.NewTicketRepository()
	scr := repository.NewScheduleRepository()
	fr := repository.NewFareRepository()
	sr := repository.NewSessionRepository()

	tc := &controller.TicketController{
		TicketUsecase: usecase.NewTicketUsecase(db, tr, scr, fr, sr),
	}
	group.POST("/ticket/create", tc.CreateTicket)
	group.GET("/tickets", tc.GetAllTickets)
	group.GET("/ticket/:id", tc.GetTicketByID)
	group.PUT("/ticket//update:id", tc.UpdateTicket)
	group.DELETE("/ticket/:id", tc.DeleteTicket)
}
