package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type PaymentRouter struct {
	Controller   *controller.PaymentController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewPaymentRouter(payment_controller *controller.PaymentController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware,
) *PaymentRouter {
	return &PaymentRouter{Controller: payment_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i PaymentRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	pc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/payment-channels", pc.GetPaymentChannels)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	protected.Use(i.Authorized.Handle())
}
