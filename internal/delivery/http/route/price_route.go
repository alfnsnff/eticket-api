package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPriceRouter(db *gorm.DB, group *gin.RouterGroup) {
	rr := repository.NewPriceRepository()
	rc := &controller.PriceController{
		PriceUsecase: usecase.NewPriceUsecase(db, rr),
	}
	group.POST("/price", rc.CreatePrice)
	group.GET("/prices", rc.GetAllPrices)
	group.GET("/price/:id", rc.GetPriceByID)
	group.GET("/price/route/:id", rc.GetPriceByRouteID)
	group.PUT("/price/:id", rc.UpdatePrice)
	group.DELETE("/price/:id", rc.DeletePrice)
}
