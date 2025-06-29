package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	requests "eticket-api/internal/delivery/http/v1/request"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	Validate       validator.Validator
	Log            logger.Logger
	BookingUsecase *usecase.BookingUsecase
}

// NewBookingController creates a new BookingController instance
func NewBookingController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	booking_usecase *usecase.BookingUsecase,

) {
	c := &BookingController{
		Log:            log,
		Validate:       validate,
		BookingUsecase: booking_usecase,
	}

	router.GET("/bookings", c.GetAllBookings)
	router.GET("/booking/:id", c.GetBookingByID)
	router.GET("/booking/order/:id", c.GetBookingByOrderID)
	router.GET("/booking/payment/callback", c.GetBookingByID)

	protected.POST("/booking/create", c.CreateBooking)
	protected.PUT("/booking/update/:id", c.UpdateBooking)
	protected.DELETE("/booking/:id", c.DeleteBooking)
}

func (c *BookingController) CreateBooking(ctx *gin.Context) {
	request := new(requests.CreateBookingRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.BookingUsecase.CreateBooking(ctx, requests.BookingFromCreate(request)); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}

		c.Log.WithError(err).Error("failed to create booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking created successfully", nil))
}

func (c *BookingController) GetAllBookings(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := c.BookingUsecase.ListBookings(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve bookings")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve bookings", err.Error()))
		return
	}

	responses := make([]*requests.BookingResponse, len(datas))
	for i, data := range datas {
		responses[i] = requests.BookingToResponse(data)
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		responses,
		"Bookings retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

// GetBookingByID retrieves a booking by its ID
func (c *BookingController) GetBookingByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid booking ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	data, err := c.BookingUsecase.GetBookingByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}
		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(requests.BookingToResponse(data), "Booking retrieved successfully", nil))
}

// GetBookingByOrderID retrieves a booking by its Order ID
func (c *BookingController) GetBookingByOrderID(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := c.BookingUsecase.GetBookingByOrderID(ctx, id)

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}

		c.Log.WithError(err).WithField("orderID", id).Error("failed to retrieve booking by order ID")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if data == nil {
		c.Log.WithField("orderID", id).Warn("booking not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(requests.BookingToResponse(data), "Booking retrieved successfully", nil))
}

func (c *BookingController) UpdateBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid or missing booking ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing booking ID", nil))
		return
	}

	request := new(requests.UpdateBookingRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.BookingUsecase.UpdateBooking(ctx, requests.BookingFromUpdate(request)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("booking already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("booking already exists", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to update booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking updated successfully", nil))
}

func (c *BookingController) DeleteBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid booking ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	if err := c.BookingUsecase.DeleteBooking(ctx, uint(id)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to delete booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking deleted successfully", nil))
}
