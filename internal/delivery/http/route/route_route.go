package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouteRouter(db *gorm.DB, group *gin.RouterGroup) {
	rr := repository.NewRouteRepository()
	rc := &controller.RouteController{
		RouteUsecase: usecase.NewRouteUsecase(db, rr),
	}
	group.POST("/route/create", rc.CreateRoute)
	group.GET("/routes", rc.GetAllRoutes)
	group.GET("/route/:id", rc.GetRouteByID)
	group.PUT("/route//update:id", rc.UpdateRoute)
	group.DELETE("/route/:id", rc.DeleteRoute)
}
