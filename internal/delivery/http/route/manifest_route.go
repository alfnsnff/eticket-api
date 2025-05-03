package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewCapacityRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewManifestRepository()
	hc := &controller.ManifestController{
		ManifestUsecase: usecase.NewManifestUsecase(db, hr),
	}
	group.POST("/manifest", hc.CreateManifest)
	group.GET("/manifests", hc.GetAllManifests)
	group.GET("/manifest/:id", hc.GetManifestByID)
	group.PUT("/manifest/:id", hc.UpdateManifest)
	group.DELETE("/manifest/:id", hc.DeleteManifest)
}
