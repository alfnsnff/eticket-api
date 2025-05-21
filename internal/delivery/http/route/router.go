package route

import (
	authrouter "eticket-api/internal/delivery/http/route/auth"
	"eticket-api/internal/injector"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, ic *injector.Container) {

	group := router.Group("/api/v1")
	NewTicketRouter(ic, group)
	NewRouteRouter(ic, group)
	NewClassRouter(ic, group)
	NewHarborRouter(ic, group)
	NewBookingRouter(ic, group)
	NewScheduleRouter(ic, group)
	NewCapacityRouter(ic, group)
	NewFareRouter(ic, group)
	NewAllocationRouter(ic, group)
	NewSessionRouter(ic, group)

	NewShipRouter(ic, group)

	authrouter.NewAuthRouter(ic, group)
	authrouter.NewRoleRouter(ic, group)
	authrouter.NewUserRouter(ic, group)
	authrouter.NewUserRoleRouter(ic, group)

}
