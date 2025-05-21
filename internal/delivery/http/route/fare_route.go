package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewFareRouter(ic *injector.Container, rg *gin.RouterGroup) {
	rr := ic.FareRepository
	fc := &controller.FareController{
		FareUsecase: usecase.NewFareUsecase(ic.FareUsecase.Tx, rr),
	}

	public := rg.Group("") // No middleware
	public.GET("/fares", fc.GetAllFares)
	public.GET("/fare/:id", fc.GetFareByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/fare/create", fc.CreateFare)
	protected.PUT("/fare/update/:id", fc.UpdateFare)
	protected.DELETE("/fare/:id", fc.DeleteFare)
}
