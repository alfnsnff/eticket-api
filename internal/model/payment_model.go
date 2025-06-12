package model

import "encoding/json"

type Result struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type Instruction struct {
	Title string   `json:"title"`
	Steps []string `json:"steps"`
}

type OrderItem struct {
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	Subtotal   int    `json:"subtotal"`
	ProductURL string `json:"product_url"`
	ImageURL   string `json:"image_url"`
}

type WriteTransactionRequest struct {
	Method        string      `json:"method"`
	MerchantRef   string      `json:"merchant_ref"`
	Amount        int         `json:"amount"`
	CustomerName  string      `json:"customer_name"`
	CustomerPhone string      `json:"customer_phone"`
	CustomerEmail string      `json:"customer_email"`
	OrderItems    []OrderItem `json:"order_items"`
	CallbackUrl   string      `json:"callback_url"`
	ReturnUrl     string      `json:"return_url"`
	ExpiredTime   int         `json:"expired_time"`
	Signature     string      `json:"signature"`
}

type ReadTransactionResponse struct {
	Reference            string        `json:"reference"`
	MerchantRef          string        `json:"merchant_ref"`
	PaymentSelectionType string        `json:"payment_selection_type"`
	PaymentMethod        string        `json:"payment_method"`
	PaymentName          string        `json:"payment_name"`
	CustomerName         string        `json:"customer_name"`
	CustomerEmail        string        `json:"customer_email"`
	CustomerPhone        string        `json:"customer_phone"`
	CallbackUrl          string        `json:"callback_url"`
	ReturnUrl            string        `json:"return_url"`
	Amount               int           `json:"amount"`
	FeeMerchant          int           `json:"fee_merchant"`
	FeeCustomer          int           `json:"fee_customer"`
	TotalFee             int           `json:"total_fee"`
	AmountReceived       int           `json:"amount_received"`
	PayCode              string        `json:"pay_code"`
	PayUrl               *string       `json:"pay_url"`
	CheckoutUrl          string        `json:"checkout_url"`
	Status               string        `json:"status"`
	ExpiredTime          int64         `json:"expired_time"`
	OrderItems           []OrderItem   `json:"order_items"`
	Instructions         []Instruction `json:"instructions"`
	QrString             *string       `json:"qr_string"`
	QrUrl                *string       `json:"qr_url"`
}

type ReadPaymentChannelResponse struct {
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

type Fee struct {
	Flat    float64 `json:"flat"`
	Percent float64 `json:"percent"`
}

type FeePercent struct {
	Flat    float64 `json:"flat"`
	Percent string  `json:"percent"` // Keep as string since JSON uses quotes
}
