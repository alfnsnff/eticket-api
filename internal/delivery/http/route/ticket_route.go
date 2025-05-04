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
	sr := repository.NewScheduleRepository()
	fr := repository.NewFareRepository()

	tc := &controller.TicketController{
		TicketUsecase: usecase.NewTicketUsecase(db, tr, sr, fr),
	}
	group.POST("/ticket", tc.CreateTicket)
	group.POST("/ticket/fill", tc.FillTicketData)
	group.GET("/tickets", tc.GetAllTickets)
	group.GET("/ticket/:id", tc.GetTicketByID)
	group.PUT("/ticket/:id", tc.UpdateTicket)
	group.DELETE("/ticket/:id", tc.DeleteTicket)
}
