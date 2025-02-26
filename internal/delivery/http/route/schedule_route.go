package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewScheduleRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewScheduleRepository(db)
	hc := &controller.ScheduleController{
		ScheduleUsecase: usecase.NewScheduleUsecase(hr),
	}
	group.POST("/schedule", hc.CreateSchedule)
	group.GET("/schedules", hc.GetAllSchedules)
	group.GET("/schedule/:id", hc.GetScheduleByID)
	group.PUT("/schedule/:id", hc.UpdateSchedule)
	group.DELETE("/schedule/:id", hc.DeleteSchedule)
}
