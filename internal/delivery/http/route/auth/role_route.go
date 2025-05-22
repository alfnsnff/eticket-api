package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	usecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

func NewRoleRouter(ic *injector.Container, rg *gin.RouterGroup) {
	ror := ic.Repository.RoleRepository
	roc := &authcontroller.RoleController{
		RoleUsecase: usecase.NewRoleUsecase(ic.Tx, ror),
	}
	public := rg.Group("") // No middleware
	public.GET("/roles", roc.GetAllRoles)
	public.GET("/role/:id", roc.GetRoleByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.Repository.UserRepository, ic.Repository.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/role/create", roc.CreateRole)
	protected.PUT("/role/update/:id", roc.UpdateRole)
	protected.DELETE("/role/:id", roc.DeleteRole)
}
