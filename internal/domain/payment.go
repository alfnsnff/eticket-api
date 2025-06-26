package domain

import "eticket-api/internal/model"

type TripayClient interface {
	CreatePayment(payload *model.WriteTransactionRequest) (model.ReadTransactionResponse, error)
	GetPaymentChannels() ([]model.ReadPaymentChannelResponse, error)
	GetTransactionDetail(reference string) (*model.ReadTransactionResponse, error)
}
