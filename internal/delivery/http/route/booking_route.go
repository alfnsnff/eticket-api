package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewBookingRouter(ic *injector.Container, rg *gin.RouterGroup) {
	br := ic.Repository.BookingRepository
	tr := ic.Repository.TicketRepository
	csr := ic.Repository.SessionRepository
	bc := &controller.BookingController{
		BookingUsecase: usecase.NewBookingUsecase(ic.Tx, br, tr, csr),
	}

	public := rg.Group("") // No middleware
	// public.POST("/booking/confirm", bc.ConfirmBooking)
	public.GET("/bookings", bc.GetAllBookings)
	public.GET("/booking/:id", bc.GetBookingByID)
	public.GET("/booking/payment/callback", bc.GetBookingByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager)
	protected.Use(middleware.Authenticate())

	protected.POST("/booking/create", bc.CreateBooking)
	protected.PUT("/booking/update/:id", bc.UpdateBooking)
	protected.DELETE("/booking/:id", bc.DeleteBooking)
}
