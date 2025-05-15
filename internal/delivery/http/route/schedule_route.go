package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewScheduleRouter(db *gorm.DB, group *gin.RouterGroup) {
	ar := repository.NewAllocationRepository()
	cr := repository.NewClassRepository()
	fr := repository.NewFareRepository()
	mr := repository.NewManifestRepository()
	rr := repository.NewRouteRepository()
	shr := repository.NewShipRepository()
	scr := repository.NewScheduleRepository()
	tr := repository.NewTicketRepository()
	hc := &controller.ScheduleController{
		ScheduleUsecase: usecase.NewScheduleUsecase(db, ar, cr, fr, mr, rr, shr, scr, tr),
	}
	group.POST("/schedule/create", hc.CreateSchedule)
	group.POST("/schedule/allocation/create", hc.CreateScheduleWithAllocation)
	group.GET("/schedules", hc.GetAllSchedules)
	group.GET("/schedules/active", hc.GetAllScheduled)
	group.GET("/schedule/:id", hc.GetScheduleByID)
	group.GET("/schedule/:id/quota", hc.GetQuotaByScheduleID)
	group.PUT("/schedule/update/:id", hc.UpdateSchedule)
	group.DELETE("/schedule/:id", hc.DeleteSchedule)
}
