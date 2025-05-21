package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewHarborRouter(ic *injector.Container, rg *gin.RouterGroup) {
	hr := ic.HarborRepository
	hc := &controller.HarborController{
		HarborUsecase: usecase.NewHarborUsecase(ic.HarborUsecase.Tx, hr),
	}

	public := rg.Group("") // No middleware
	public.GET("/harbors", hc.GetAllHarbors)
	public.GET("/harbor/:id", hc.GetHarborByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/harbor/create", hc.CreateHarbor)
	protected.PUT("/harbor/update/:id", hc.UpdateHarbor)
	protected.DELETE("/harbor/:id", hc.DeleteHarbor)
}
