package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	Validate       validator.Validator
	Log            logger.Logger
	PaymentUsecase *usecase.PaymentUsecase
}

// NewPaymentController creates a new PaymentController instance
func NewPaymentController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	payment_usecase *usecase.PaymentUsecase,

) {
	c := &PaymentController{
		Log:            log,
		Validate:       validate,
		PaymentUsecase: payment_usecase,
	}

	router.GET("/payment-channels", c.GetPaymentChannels)
	router.GET("/payment/transaction/detail/:id", c.GetTransactionDetail)
	router.POST("/payment/transaction/create", c.CreatePayment)
	router.POST("/payment/callback", c.HandleCallback)
}

func (c *PaymentController) GetPaymentChannels(ctx *gin.Context) {
	datas, err := c.PaymentUsecase.ListPaymentChannels(ctx)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve payment channels")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve payment channels", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Payment channels retrieved successfully", nil))
}

func (c *PaymentController) GetTransactionDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	datas, err := c.PaymentUsecase.GetTransactionDetail(ctx, id)

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("transaction not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("transaction not found", nil))
			return
		}

		c.Log.WithError(err).WithField("reference", id).Error("failed to retrieve transaction detail")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve transaction detail", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Transaction detail retrieved successfully", nil))
}

func (c *PaymentController) CreatePayment(ctx *gin.Context) {

	request := new(model.WritePaymentRequest)
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

	datas, err := c.PaymentUsecase.CreatePayment(ctx, request)
	if err != nil {
		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}
		c.Log.WithError(err).Error("failed to create payment transaction")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to initiate transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Transaction initiated successfully", nil))
}

func (c *PaymentController) HandleCallback(ctx *gin.Context) {
	request := new(model.WriteCallbackRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON callback request")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	err := c.PaymentUsecase.HandleCallback(ctx, request)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", request.MerchantRef).Warn("payment not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("payment not found", nil))
			return
		}

		c.Log.WithError(err).WithField("reference", request.Reference).Error("failed to handle payment callback")
		ctx.JSON(http.StatusForbidden, response.NewErrorResponse("Callback handling failed", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Callback verified", nil))
}
