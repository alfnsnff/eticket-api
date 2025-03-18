package controller

import (
	"eticket-api/internal/domain/dto"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response" // Import the response package
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	BookingUsecase usecase.BookingUsecase
}

// CreateBooking handles creating a new Booking
func (h *BookingController) CreateBooking(c *gin.Context) {
	var booking entities.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := h.BookingUsecase.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking created successfully", nil))
}

// GetAllBookings handles retrieving all Bookings
func (h *BookingController) GetAllBookings(c *gin.Context) {
	bookings, err := h.BookingUsecase.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve bookings", err.Error()))
		return
	}

	bookingDTOs := dto.ToBookingDTOs(bookings)
	c.JSON(http.StatusOK, response.NewSuccessResponse(bookingDTOs, "Bookings retrieved successfully", nil))
}

// GetBookingByID handles retrieving a Booking by its ID
func (h *BookingController) GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	booking, err := h.BookingUsecase.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if booking == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	bookingDTO := dto.ToBookingDTO(booking)
	c.JSON(http.StatusOK, response.NewSuccessResponse(bookingDTO, "Booking retrieved successfully", nil))
}

// UpdateBooking handles updating an existing Booking
func (h *BookingController) UpdateBooking(c *gin.Context) {
	var booking entities.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if booking.ID == 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Booking ID is required", nil))
		return
	}

	if err := h.BookingUsecase.UpdateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update booking", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking updated successfully", nil))
}

// DeleteBooking handles deleting a Booking by its ID
func (h *BookingController) DeleteBooking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	if err := h.BookingUsecase.DeleteBooking(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete booking", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking deleted successfully", nil))
}

// NewBookingController creates a new BookingController instance.
func NewBookingController(bookingUsecase usecase.BookingUsecase) *BookingController {
	return &BookingController{BookingUsecase: bookingUsecase}
}
