package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	Controller   *controller.PaymentController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewPaymentRoute(payment_controller *controller.PaymentController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware,
) *PaymentRoute {
	return &PaymentRoute{Controller: payment_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i PaymentRoute) Set(router *gin.Engine) {
	pc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/payment-channels", pc.GetPaymentChannels)
	public.GET("/payment/transaction/detail/:id", pc.GetTransactionDetail)
	public.POST("/payment/transaction/create", pc.CreatePayment)
	public.POST("/payment/callback", pc.HandleCallback)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	protected.Use(i.Authorized.Handle())
}
