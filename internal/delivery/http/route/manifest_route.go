package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewCapacityRouter(ic *injector.Container, rg *gin.RouterGroup) {
	hr := ic.ManifestRepository
	mc := &controller.ManifestController{
		ManifestUsecase: usecase.NewManifestUsecase(ic.ManifestUsecase.Tx, hr),
	}

	public := rg.Group("") // No middleware
	public.GET("/manifests", mc.GetAllManifests)
	public.GET("/manifest/:id", mc.GetManifestByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/manifest/create", mc.CreateManifest)
	protected.PUT("/manifest/update/:id", mc.UpdateManifest)
	protected.DELETE("/manifest/:id", mc.DeleteManifest)
}
