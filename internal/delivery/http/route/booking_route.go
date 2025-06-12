package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type BookingRouter struct {
	Controller   *controller.BookingController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewBookingRouter(booking_controller *controller.BookingController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *BookingRouter {
	return &BookingRouter{Controller: booking_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i BookingRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {

	public := rg.Group("") // No middleware
	// public.POST("/booking/confirm", bc.ConfirmBooking)
	public.GET("/bookings", i.Controller.GetAllBookings)
	public.GET("/booking/:id", i.Controller.GetBookingByID)
	public.GET("/booking/payment/callback", i.Controller.GetBookingByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/booking/create", i.Controller.CreateBooking)
	protected.PUT("/booking/update/:id", i.Controller.UpdateBooking)
	protected.DELETE("/booking/:id", i.Controller.DeleteBooking)
}
