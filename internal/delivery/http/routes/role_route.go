package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RoleRoute struct {
	Controller   *controller.RoleController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewRoleRoute(role_controller *controller.RoleController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware,
) *RoleRoute {
	return &RoleRoute{Controller: role_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i RoleRoute) Set(router *gin.Engine) {
	roc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/roles", roc.GetAllRoles)
	public.GET("/role/:id", roc.GetRoleByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	protected.Use(i.Authorized.Handle())

	protected.POST("/role/create", roc.CreateRole)
	protected.PUT("/role/update/:id", roc.UpdateRole)
	protected.DELETE("/role/:id", roc.DeleteRole)
}
