package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/usecase/payment"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentUsecase *payment.PaymentUsecase
}

// NewPaymentController creates a new PaymentController instance
func NewPaymentController(payment_usecase *payment.PaymentUsecase) *PaymentController {
	return &PaymentController{PaymentUsecase: payment_usecase}
}

func (pc *PaymentController) GetPaymentChannels(ctx *gin.Context) {
	datas, err := pc.PaymentUsecase.GetPaymentChannels()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve roles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Roles retrieved successfully", err))
}
