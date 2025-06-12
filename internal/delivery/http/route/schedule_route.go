package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ScheduleRouter struct {
	Controller   *controller.ScheduleController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewScheduleRouter(schedule_controller *controller.ScheduleController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ScheduleRouter {
	return &ScheduleRouter{Controller: schedule_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ScheduleRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	scc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/schedules", scc.GetAllSchedules)
	public.GET("/schedules/active", scc.GetAllScheduled)
	public.GET("/schedule/:id", scc.GetScheduleByID)
	public.GET("/schedule/:id/quota", scc.GetQuotaByScheduleID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/schedule/create", scc.CreateSchedule)
	protected.POST("/schedule/Schedule/create", scc.CreateScheduleWithAllocation)
	protected.PUT("/schedule/update/:id", scc.UpdateSchedule)
	protected.DELETE("/schedule/:id", scc.DeleteSchedule)
}
