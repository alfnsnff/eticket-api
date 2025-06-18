package controller

import (
	"eticket-api/internal/common/response"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase/payment"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentUsecase *payment.PaymentUsecase
	Authenticate   *middleware.AuthenticateMiddleware
	Authorized     *middleware.AuthorizeMiddleware
}

// NewPaymentController creates a new PaymentController instance
func NewPaymentController(
	g *gin.Engine, payment_usecase *payment.PaymentUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	pc := &PaymentController{
		PaymentUsecase: payment_usecase,
		Authenticate:   authtenticate,
		Authorized:     authorized}

	public := g.Group("") // No middleware
	public.GET("/payment-channels", pc.GetPaymentChannels)
	public.GET("/payment/transaction/detail/:id", pc.GetTransactionDetail)
	public.POST("/payment/transaction/create", pc.CreatePayment)
	public.POST("/payment/callback", pc.HandleCallback)

	protected := g.Group("")
	protected.Use(pc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())
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
	request := new(model.WritePaymentRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	datas, err := pc.PaymentUsecase.CreatePayment(ctx, request)

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
