package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTicketRouter(db *gorm.DB, group *gin.RouterGroup) {
	tr := repository.NewTicketRepository(db)
	sr := repository.NewScheduleRepository(db)
	scr := repository.NewShipClassRepository(db)
	pr := repository.NewPriceRepository(db)

	tc := &controller.TicketController{
		TicketUsecase: usecase.NewTicketUsecase(tr, sr, scr, pr),
	}
	group.POST("/ticket", tc.CreateTicket)
	group.POST("/ticket/validate", tc.ValidateTicket)
	group.POST("/ticket/check", tc.GetBookedCount)
	group.GET("/tickets", tc.GetAllTickets)
	group.GET("/ticket/:id", tc.GetTicketByID)
	group.PUT("/ticket/:id", tc.UpdateTicket)
	group.DELETE("/ticket/:id", tc.DeleteTicket)
}
