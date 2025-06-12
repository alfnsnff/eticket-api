package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RoleRouter struct {
	Controller   *controller.RoleController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewRoleRouter(role_controller *controller.RoleController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware,
) *RoleRouter {
	return &RoleRouter{Controller: role_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i RoleRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	roc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/roles", roc.GetAllRoles)
	public.GET("/role/:id", roc.GetRoleByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	protected.Use(i.Authorized.Handle())

	protected.POST("/role/create", roc.CreateRole)
	protected.PUT("/role/update/:id", roc.UpdateRole)
	protected.DELETE("/role/:id", roc.DeleteRole)
}
