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
	tc := &controller.TicketController{
		TicketUsecase: usecase.NewTicketUsecase(tr),
	}
	group.POST("/tickets", tc.CreateTicket)
	group.GET("/tickets", tc.GetAllTickets)
	group.GET("/tickets/:id", tc.GetTicketByID)
	group.PUT("/tickets/:id", tc.UpdateTicket)
	group.DELETE("/tickets/:id", tc.DeleteTicket)
}
