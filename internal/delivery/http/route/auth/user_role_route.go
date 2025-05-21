package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	usecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

func NewUserRoleRouter(ic *injector.Container, rg *gin.RouterGroup) {
	rr := ic.RoleRepository
	ur := ic.UserRepository
	urr := ic.UserRoleRepository
	urc := &authcontroller.UserRoleController{
		UserRoleUsecase: usecase.NewUserRoleUsecase(ic.UserRoleUsecase.Tx, rr, ur, urr),
	}

	public := rg.Group("") // No middleware
	public.GET("/user/roles", urc.GetAllUserRoles)
	public.GET("/user/role/:id", urc.GetUserRoleByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/user/role/assign", urc.CreateUserRole)
	protected.PUT("/user/role/update/:id", urc.UpdateUserRole)
	protected.DELETE("/user/role/:id", urc.DeleteUserRole)
}
