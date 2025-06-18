package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type BookingRoute struct {
	Controller   *controller.BookingController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewBookingRoute(booking_controller *controller.BookingController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *BookingRoute {
	return &BookingRoute{Controller: booking_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i BookingRoute) Set(router *gin.Engine) {

	public := router.Group("") // No middleware
	// public.POST("/booking/confirm", bc.ConfirmBooking)
	public.GET("/bookings", i.Controller.GetAllBookings)
	public.GET("/booking/:id", i.Controller.GetBookingByID)
	public.GET("/booking/order/:id", i.Controller.GetBookingByOrderID)
	public.GET("/booking/payment/callback", i.Controller.GetBookingByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/booking/create", i.Controller.CreateBooking)
	protected.PUT("/booking/update/:id", i.Controller.UpdateBooking)
	protected.DELETE("/booking/:id", i.Controller.DeleteBooking)
}
