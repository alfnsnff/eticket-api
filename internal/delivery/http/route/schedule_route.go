package route

import (
	"eticket-api/internal/delivery/http/controller"
	// "eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewScheduleRouter(db *gorm.DB, group *gin.RouterGroup) {
	// hs := repository.NewScheduleRepository(db)
	// hp := repository.NewPriceRepository(db)
	// ht := repository.NewTicketRepository(db)
	hc := &controller.ScheduleController{
		ScheduleUsecase: usecase.ScheduleUsecase{},
	}
	group.POST("/schedule", hc.CreateSchedule)
	group.POST("/schedule/search", hc.SearchSchedule)
	group.GET("/schedule/quota/schedule/:id", hc.GetPricesWithQuota)
	group.GET("/schedules", hc.GetAllSchedules)
	group.GET("/schedule/:id", hc.GetScheduleByID)
	group.PUT("/schedule/:id", hc.UpdateSchedule)
	group.DELETE("/schedule/:id", hc.DeleteSchedule)
}
