package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewBookingRouter(ic *injector.Container, rg *gin.RouterGroup) {
	br := ic.BookingRepository
	tr := ic.TicketRepository
	csr := ic.SessionRepository
	bc := &controller.BookingController{
		BookingUsecase: usecase.NewBookingUsecase(ic.BookingUsecase.Tx, br, tr, csr),
	}

	public := rg.Group("") // No middleware
	public.POST("/booking/confirm", bc.ConfirmBooking)
	public.GET("/bookings", bc.GetAllBookings)
	public.GET("/booking/:id", bc.GetBookingByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/booking/create", bc.CreateBooking)
	protected.PUT("/booking/update/:id", bc.UpdateBooking)
	protected.DELETE("/booking/:id", bc.DeleteBooking)
}
