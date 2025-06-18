package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ScheduleRoute struct {
	Controller   *controller.ScheduleController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewScheduleRoute(schedule_controller *controller.ScheduleController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ScheduleRoute {
	return &ScheduleRoute{Controller: schedule_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ScheduleRoute) Set(router *gin.Engine) {
	scc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/schedules", scc.GetAllSchedules)
	public.GET("/schedules/active", scc.GetAllScheduled)
	public.GET("/schedule/:id", scc.GetScheduleByID)
	public.GET("/schedule/:id/quota", scc.GetQuotaByScheduleID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/schedule/create", scc.CreateSchedule)
	protected.POST("/schedule/Schedule/create", scc.CreateScheduleWithAllocation)
	protected.PUT("/schedule/update/:id", scc.UpdateSchedule)
	protected.DELETE("/schedule/:id", scc.DeleteSchedule)
}
