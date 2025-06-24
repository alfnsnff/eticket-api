package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/payment"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	Log            logger.Logger
	Validate       validator.Validator
	PaymentUsecase *payment.PaymentUsecase
}

// NewPaymentController creates a new PaymentController instance
func NewPaymentController(
	log logger.Logger,
	validate validator.Validator,
	payment_usecase *payment.PaymentUsecase,
) *PaymentController {
	return &PaymentController{
		Log:            log,
		Validate:       validate,
		PaymentUsecase: payment_usecase,
	}

}

func (pc *PaymentController) GetPaymentChannels(ctx *gin.Context) {
	datas, err := pc.PaymentUsecase.GetPaymentChannels(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve roles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Roles retrieved successfully", err))
}

func (pc *PaymentController) GetTransactionDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	datas, err := pc.PaymentUsecase.GetTransactionDetail(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve roles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Roles retrieved successfully", err))
}

func (pc *PaymentController) CreatePayment(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing session id", err.Error()))
		return
	}

	request := new(model.WritePaymentRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := pc.Validate.Struct(request); err != nil {
		pc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, err := pc.PaymentUsecase.CreatePayment(ctx, request, sessionID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to initiate transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Transaction initiate successfully", err))
}

func (pc *PaymentController) HandleCallback(ctx *gin.Context) {
	request := new(model.WriteCallbackRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}
	err := pc.PaymentUsecase.HandleCallback(ctx, ctx.Request, request)
	if err != nil {
		ctx.JSON(http.StatusForbidden, response.NewErrorResponse("Callback handling failed", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Callback verified", nil))
}
