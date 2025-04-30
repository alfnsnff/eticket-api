package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewFareRouter(db *gorm.DB, group *gin.RouterGroup) {
	rr := repository.NewFareRepository()
	rc := &controller.FareController{
		FareUsecase: usecase.NewFareUsecase(db, rr),
	}
	group.POST("/fare", rc.CreateFare)
	group.GET("/fares", rc.GetAllFares)
	group.GET("/fare/:id", rc.GetFareByID)
	group.GET("/fares/route/:id", rc.GetFareByRouteID)
	group.PUT("/fare/:id", rc.UpdateFare)
	group.DELETE("/fare/:id", rc.DeleteFare)
}
