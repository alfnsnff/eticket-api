package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ManifestRouter struct {
	Controller   *controller.ManifestController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewManifestRouter(manifest_controller *controller.ManifestController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ManifestRouter {
	return &ManifestRouter{Controller: manifest_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ManifestRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	mc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/manifests", mc.GetAllManifests)
	public.GET("/manifest/:id", mc.GetManifestByID)

	protected := rg.Group("")
	protected.Use(i.Authorized.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/manifest/create", mc.CreateManifest)
	protected.PUT("/manifest/update/:id", mc.UpdateManifest)
	protected.DELETE("/manifest/:id", mc.DeleteManifest)
}
