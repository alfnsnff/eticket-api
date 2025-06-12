package payment

import (
	"eticket-api/internal/client"
	"eticket-api/internal/model"
)

type PaymentUsecase struct {
	TripayClient *client.TripayClient
}

func NewPaymentUsecase(
	tripay_client *client.TripayClient,
) *PaymentUsecase {
	return &PaymentUsecase{
		TripayClient: tripay_client,
	}
}

func (c *PaymentUsecase) GetPaymentChannels() ([]*model.ReadPaymentChannelResponse, error) {
	channels, err := c.TripayClient.GetPaymentChannels()
	if err != nil {
		return nil, err
	}
	result := make([]*model.ReadPaymentChannelResponse, len(channels))
	for i := range channels {
		result[i] = &channels[i]
	}
	return result, nil
}
