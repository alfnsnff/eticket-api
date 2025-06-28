package request

// import "encoding/json"

// type Result struct {
// 	Success bool            `json:"success"`
// 	Message string          `json:"message"`
// 	Data    json.RawMessage `json:"data"`
// }

// type Instruction struct {
// 	Title string   `json:"title"`
// 	Steps []string `json:"steps"`
// }

// type Fee struct {
// 	Flat    float64 `json:"flat"`
// 	Percent float64 `json:"percent"`
// }

// type FeePercent struct {
// 	Flat    float64 `json:"flat"`
// 	Percent string  `json:"percent"` // Keep as string since JSON uses quotes
// }

// type ReadPaymentChannelResponse struct {
// 	Group         string     `json:"group"`
// 	Code          string     `json:"code"`
// 	Name          string     `json:"name"`
// 	Type          string     `json:"type"`
// 	FeeMerchant   Fee        `json:"fee_merchant"`
// 	FeeCustomer   Fee        `json:"fee_customer"`
// 	TotalFee      FeePercent `json:"total_fee"`
// 	MinimumFee    float64    `json:"minimum_fee"`
// 	MaximumFee    float64    `json:"maximum_fee"`
// 	MinimumAmount float64    `json:"minimum_amount"`
// 	MaximumAmount float64    `json:"maximum_amount"`
// 	IconURL       string     `json:"icon_url"`
// 	Active        bool       `json:"active"`
// }

// type WritePaymentRequest struct {
// 	OrderID       string `json:"order_id"`
// 	PaymentMethod string `json:"payment_method"`
// }

// type OrderItem struct {
// 	SKU        string `json:"sku" validate:"omitempty"`
// 	Name       string `json:"name" validate:"required"`
// 	Price      int    `json:"price" validate:"required,gt=0"`
// 	Quantity   int    `json:"quantity" validate:"required,gt=0"`
// 	Subtotal   int    `json:"subtotal" validate:"required,gt=0"`
// 	ProductURL string `json:"product_url" validate:"omitempty,url"`
// 	ImageURL   string `json:"image_url" validate:"omitempty,url"`
// }

// type WriteTransactionRequest struct {
// 	Method        string      `json:"method" validate:"required"`
// 	MerchantRef   string      `json:"merchant_ref" validate:"required"`
// 	Amount        int         `json:"amount" validate:"required,gt=0"`
// 	CustomerName  string      `json:"customer_name" validate:"required"`
// 	CustomerPhone string      `json:"customer_phone" validate:"omitempty,e164"` // Optional but valid if provided
// 	CustomerEmail string      `json:"customer_email" validate:"required,email"`
// 	OrderItems    []OrderItem `json:"order_items" validate:"required,dive"`
// 	CallbackUrl   string      `json:"callback_url" validate:"omitempty,url"`
// 	ReturnUrl     string      `json:"return_url" validate:"omitempty,url"`
// 	ExpiredTime   int         `json:"expired_time" validate:"omitempty,gt=0"`
// 	Signature     string      `json:"signature" validate:"required"`
// }

// type ReadTransactionResponse struct {
// 	Reference            string        `json:"reference"`
// 	MerchantRef          string        `json:"merchant_ref"`
// 	PaymentSelectionType string        `json:"payment_selection_type"`
// 	PaymentMethod        string        `json:"payment_method"`
// 	PaymentName          string        `json:"payment_name"`
// 	CustomerName         string        `json:"customer_name"`
// 	CustomerEmail        string        `json:"customer_email"`
// 	CustomerPhone        string        `json:"customer_phone"`
// 	CallbackUrl          string        `json:"callback_url"`
// 	ReturnUrl            string        `json:"return_url"`
// 	Amount               int           `json:"amount"`
// 	FeeMerchant          int           `json:"fee_merchant"`
// 	FeeCustomer          int           `json:"fee_customer"`
// 	TotalFee             int           `json:"total_fee"`
// 	AmountReceived       int           `json:"amount_received"`
// 	PayCode              string        `json:"pay_code"`
// 	PayUrl               *string       `json:"pay_url"`
// 	CheckoutUrl          string        `json:"checkout_url"`
// 	Status               string        `json:"status"`
// 	ExpiredTime          int64         `json:"expired_time"`
// 	OrderItems           []OrderItem   `json:"order_items"`
// 	Instructions         []Instruction `json:"instructions"`
// 	QrString             *string       `json:"qr_string"`
// 	QrUrl                *string       `json:"qr_url"`
// }

// // TripayCallbackHandler.go
// type WriteCallbackRequest struct {
// 	Reference     string `json:"reference"`
// 	MerchantRef   string `json:"merchant_ref"`
// 	Status        string `json:"status"`
// 	Amount        int    `json:"amount"`
// 	PaymentMethod string `json:"payment_method"`
// 	Signature     string `json:"signature"`
// }

// // TripayCallbackHandler.go
// type ReadCallbackResponse struct {
// 	Success bool `json:"success"`
// }
