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
	pc := &PaymentController{
		Log:            log,
		Validate:       validate,
		PaymentUsecase: payment_usecase,
	}

	router.GET("/payment-channels", pc.GetPaymentChannels)
	router.GET("/payment/transaction/detail/:id", pc.GetTransactionDetail)
	router.POST("/payment/transaction/create", pc.CreatePayment)
	router.POST("/payment/transaction/create/test", pc.TESTCreatePayment)
	router.POST("/payment/callback", pc.HandleCallback)
}

func (pc *PaymentController) GetPaymentChannels(ctx *gin.Context) {
	datas, err := pc.PaymentUsecase.ListPaymentChannels(ctx)

	if err != nil {
		pc.Log.WithError(err).Error("failed to retrieve payment channels")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve payment channels", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Payment channels retrieved successfully", nil))
}

func (pc *PaymentController) GetTransactionDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	datas, err := pc.PaymentUsecase.GetTransactionDetail(ctx, id)

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			pc.Log.WithField("id", id).Warn("transaction not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("transaction not found", nil))
			return
		}

		pc.Log.WithError(err).WithField("reference", id).Error("failed to retrieve transaction detail")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve transaction detail", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Transaction detail retrieved successfully", nil))
}

func (pc *PaymentController) TESTCreatePayment(ctx *gin.Context) {
	orderID, err := ctx.Cookie("order_id")
	if err != nil {
		pc.Log.WithError(err).Error("missing order ID in request")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing order id", err.Error()))
		return
	}

	request := new(model.WritePaymentRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		pc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := pc.Validate.Struct(request); err != nil {
		pc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, err := pc.PaymentUsecase.TESTCreatePayment(ctx, request, orderID)
	if err != nil {
		pc.Log.WithError(err).WithField("sessionID", orderID).Error("failed to create payment transaction")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to initiate transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Transaction initiated successfully", nil))
}

func (pc *PaymentController) CreatePayment(ctx *gin.Context) {

	request := new(model.WritePaymentRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		pc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := pc.Validate.Struct(request); err != nil {
		pc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, err := pc.PaymentUsecase.CreatePayment(ctx, request)
	if err != nil {
		pc.Log.WithError(err).Error("failed to create payment transaction")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to initiate transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Transaction initiated successfully", nil))
}

func (pc *PaymentController) HandleCallback(ctx *gin.Context) {
	request := new(model.WriteCallbackRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		pc.Log.WithError(err).Error("failed to bind JSON callback request")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	err := pc.PaymentUsecase.HandleCallback(ctx, request)
	if err != nil {
		pc.Log.WithError(err).WithField("reference", request.Reference).Error("failed to handle payment callback")
		ctx.JSON(http.StatusForbidden, response.NewErrorResponse("Callback handling failed", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Callback verified", nil))
}
