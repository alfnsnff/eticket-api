package module

import (
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/casbin/casbin/v2"
)

type MiddlewareModule struct {
	Authenticate *middleware.AuthenticateMiddleware
	Authorize    *middleware.AuthorizeMiddleware
}

func NewMiddlewareModule(t *jwt.TokenUtil, e *casbin.Enforcer) *MiddlewareModule {
	return &MiddlewareModule{
		Authenticate: middleware.NewAuthMiddleware(t),
		Authorize:    middleware.NewAuthorizeMiddleware(e),
	}
}
