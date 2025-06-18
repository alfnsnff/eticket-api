package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ManifestRoute struct {
	Controller   *controller.ManifestController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewManifestRoute(manifest_controller *controller.ManifestController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ManifestRoute {
	return &ManifestRoute{Controller: manifest_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ManifestRoute) Set(router *gin.Engine) {
	mc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/manifests", mc.GetAllManifests)
	public.GET("/manifest/:id", mc.GetManifestByID)

	protected := router.Group("")
	protected.Use(i.Authorized.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/manifest/create", mc.CreateManifest)
	protected.PUT("/manifest/update/:id", mc.UpdateManifest)
	protected.DELETE("/manifest/:id", mc.DeleteManifest)
}
