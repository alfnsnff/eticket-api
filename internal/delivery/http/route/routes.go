package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(router *gin.Engine, db *gorm.DB) {
	group := router.Group("/api/v1")
	NewTicketRouter(db, group)
	NewRouteRouter(db, group)
	NewClassRouter(db, group)
	NewHarborRouter(db, group)
	NewBookingRouter(db, group)
	NewScheduleRouter(db, group)
	NewShipRouter(db, group)
}
