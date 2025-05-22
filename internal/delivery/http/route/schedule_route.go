package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewScheduleRouter(ic *injector.Container, rg *gin.RouterGroup) {
	ar := ic.Repository.AllocationRepository
	cr := ic.Repository.ClassRepository
	fr := ic.Repository.FareRepository
	mr := ic.Repository.ManifestRepository
	rr := ic.Repository.RouteRepository
	shr := ic.Repository.ShipRepository
	scr := ic.Repository.ScheduleRepository
	tr := ic.Repository.TicketRepository
	scc := &controller.ScheduleController{
		ScheduleUsecase: usecase.NewScheduleUsecase(ic.Tx, ar, cr, fr, mr, rr, shr, scr, tr),
	}

	public := rg.Group("") // No middleware
	public.GET("/schedules", scc.GetAllSchedules)
	public.GET("/schedules/active", scc.GetAllScheduled)
	public.GET("/schedule/:id", scc.GetScheduleByID)
	public.GET("/schedule/:id/quota", scc.GetQuotaByScheduleID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.Repository.UserRepository, ic.Repository.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/schedule/create", scc.CreateSchedule)
	protected.POST("/schedule/allocation/create", scc.CreateScheduleWithAllocation)
	protected.PUT("/schedule/update/:id", scc.UpdateSchedule)
	protected.DELETE("/schedule/:id", scc.DeleteSchedule)
}
