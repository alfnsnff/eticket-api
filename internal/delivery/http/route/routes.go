package route

import (
	"eticket-api/internal/delivery/http/middleware"
	authrouter "eticket-api/internal/delivery/http/route/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(router *gin.Engine, db *gorm.DB) {
	publicRouter := router.Group("/api/v1")
	NewTicketRouter(db, publicRouter)
	NewRouteRouter(db, publicRouter)
	NewClassRouter(db, publicRouter)
	NewHarborRouter(db, publicRouter)
	NewBookingRouter(db, publicRouter)
	NewScheduleRouter(db, publicRouter)
	NewCapacityRouter(db, publicRouter)
	NewFareRouter(db, publicRouter)
	NewAllocationRouter(db, publicRouter)
	NewSessionRouter(db, publicRouter)

	protectedRouter := router.Group("/api/v1/")
	protectedRouter.Use(middleware.Authenticate())
	NewShipRouter(db, protectedRouter)

	authrouter.NewRoleRouter(db, publicRouter)
	authrouter.NewUserRouter(db, publicRouter)
	authrouter.NewUserRoleRouter(db, publicRouter)

}
