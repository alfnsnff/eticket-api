package requests

import "eticket-api/internal/domain"

type Instruction struct {
	Title string   `json:"title"`
	Steps []string `json:"steps"`
}

type Fee struct {
	Flat    float64 `json:"flat"`
	Percent float64 `json:"percent"`
}

type FeePercent struct {
	Flat    float64 `json:"flat"`
	Percent string  `json:"percent"` // Keep as string since JSON uses quotes
}

type PaymentChannel struct {
	Group         string     `json:"group"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	FeeMerchant   Fee        `json:"fee_merchant"`
	FeeCustomer   Fee        `json:"fee_customer"`
	TotalFee      FeePercent `json:"total_fee"`
	MinimumFee    float64    `json:"minimum_fee"`
	MaximumFee    float64    `json:"maximum_fee"`
	MinimumAmount float64    `json:"minimum_amount"`
	MaximumAmount float64    `json:"maximum_amount"`
	IconURL       string     `json:"icon_url"`
	Active        bool       `json:"active"`
}

func PaymentToResponse(channel *domain.PaymentChannel) *PaymentChannel {
	return &PaymentChannel{
		Group: channel.Group,
		Code:  channel.Code,
		Name:  channel.Name,
		Type:  channel.Type,
		FeeMerchant: Fee{
			Flat:    channel.FeeMerchant.Flat,
			Percent: channel.FeeMerchant.Percent,
		},
		FeeCustomer: Fee{
			Flat:    channel.FeeCustomer.Flat,
			Percent: channel.FeeCustomer.Percent},
		TotalFee: FeePercent{
			Flat:    channel.TotalFee.Flat,
			Percent: channel.TotalFee.Percent,
		},
		MinimumFee:    channel.MinimumFee,
		MaximumFee:    channel.MaximumFee,
		MinimumAmount: channel.MinimumAmount,
		MaximumAmount: channel.MaximumAmount,
		IconURL:       channel.IconURL,
		Active:        channel.Active,
	}
}

type CreareCallbackRequest struct {
	Reference     string `json:"reference"`
	MerchantRef   string `json:"merchant_ref"`
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Signature     string `json:"signature"`
}

func CallbackFromCreate(request *CreareCallbackRequest) *domain.Callback {
	return &domain.Callback{
		Reference:     request.Reference,
		MerchantRef:   request.MerchantRef,
		Status:        request.Status,
		Amount:        request.Amount,
		PaymentMethod: request.PaymentMethod,
		Signature:     request.Signature,
	}
}
