package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTicketRouter(timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	tr := repository.NewTicketRepository(db)
	tc := &controller.TicketController{
		TicketUsecase: usecase.NewTicketUsecase(tr),
	}
	group.GET("/task", tc.GetAllTickets)
}
