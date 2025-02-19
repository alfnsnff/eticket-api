package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouteRouter(db *gorm.DB, group *gin.RouterGroup) {
	rr := repository.NewRouteRepository(db)
	rc := &controller.RouteController{
		RouteUsecase: usecase.NewRouteUsecase(rr),
	}
	group.POST("/routes", rc.CreateRoute)
	group.GET("/routes", rc.GetAllRoutes)
	group.GET("/routes/:id", rc.GetRouteByID)
	group.PUT("/routes/:id", rc.UpdateRoute)
	group.DELETE("/routes/:id", rc.DeleteRoute)
}
