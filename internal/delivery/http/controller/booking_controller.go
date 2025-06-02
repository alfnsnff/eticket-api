package controller

import (
	"encoding/json"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/helper/meta"
	"eticket-api/pkg/utils/helper/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	BookingUsecase *usecase.BookingUsecase
}

// NewBookingController creates a new BookingController instance
func NewBookingController(booking_usecase *usecase.BookingUsecase) *BookingController {
	return &BookingController{BookingUsecase: booking_usecase}
}

func (bc *BookingController) CreateBooking(ctx *gin.Context) {
	request := new(model.WriteBookingRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	err := bc.BookingUsecase.CreateBooking(ctx, request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking created successfully", nil))
}

func (bc *BookingController) GetAllBookings(ctx *gin.Context) {
	params := meta.GetParams(ctx)
	datas, total, err := bc.BookingUsecase.GetAllBookings(ctx, params.Limit, params.Offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve bookings", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewPaginatedResponse(datas, "Bookings retrieved successfully", total, params.Limit, params.Page))
}

// GetBookingByID retrieves a booking by its ID
func (bc *BookingController) GetBookingByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	data, err := bc.BookingUsecase.GetBookingByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Booking retrieved successfully", nil))
}

func (bc *BookingController) UpdateBooking(ctx *gin.Context) {
	request := new(model.UpdateBookingRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Booking ID is required", nil))
		return
	}

	err := bc.BookingUsecase.UpdateBooking(ctx, uint(id), request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking updated successfully", nil))
}

func (bc *BookingController) DeleteBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	if err := bc.BookingUsecase.DeleteBooking(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking deleted successfully", nil))
}

// func (bc *BookingController) ConfirmBooking(ctx *gin.Context) {
// 	sessionID, err := ctx.Cookie("session_id")
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing session id", err.Error()))
// 		return
// 	}

// 	// if err := ctx.ShouldBindJSON(request); err != nil {
// 	// 	ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 	// 	return
// 	// }

// 	datas, err := bc.BookingUsecase.ConfirmBooking(ctx, sessionID)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Booking confirmed successfully", nil))
// }

func (bc *BookingController) PaidBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}
	// if err := ctx.ShouldBindJSON(request); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
	// 	return
	// }

	err = bc.BookingUsecase.PaidConfirm(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking confirmed successfully", nil))
}

func (h *BookingController) HandleCallback(ctx *gin.Context, r *http.Request) {
	var callback struct {
		ID         string `json:"id"`
		Status     string `json:"status"`      // e.g., "COMPLETED"
		ExternalID string `json:"external_id"` // Your order ID (e.g., "order-123")
		Amount     int    `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Log or verify the callback
	fmt.Printf("Received QRIS callback: %+v\n", callback)

	if callback.Status == "COMPLETED" {
		// ✅ Mark booking as paid in your DB
		externalIDUint, err := strconv.ParseUint(callback.ExternalID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid external_id"})
			return
		}
		err = h.BookingUsecase.PaidConfirm(ctx, uint(externalIDUint))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
			return
		}
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking paid successfully", nil))
}
