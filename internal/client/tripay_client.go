package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"eticket-api/config"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"
	"net/http"
	"time"
)

const (
	TripayBaseURL = "https://tripay.co.id/api-sandbox"
)

type TripayClient struct {
	HTTPClient *http.Client
	Tripay     *config.Tripay
}

func NewTripayClient(httpClient *http.Client, tripay *config.Tripay) *TripayClient {
	return &TripayClient{
		HTTPClient: httpClient,
		Tripay:     tripay,
	}
}

func GenerateTransactionSignature(merchantCode, merchantRef string, amount int, privateKey string) string {
	message := merchantCode + merchantRef + fmt.Sprintf("%d", amount)
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *TripayClient) CreatePayment(method string, amount int, name string, email string, phone string, orderID string, items []*entity.Ticket) (model.ReadTransactionResponse, error) {
	payload := model.WriteTransactionRequest{
		Method:        method,
		MerchantRef:   orderID,
		Amount:        amount,
		CustomerName:  name,
		CustomerEmail: email,
		CustomerPhone: phone,
		OrderItems:    TicketsToItemsTr(items),
		CallbackUrl:   "https://example.com/callback",
		ReturnUrl:     "https://example.com/success",
		ExpiredTime:   int(time.Now().Add(30 * time.Minute).Unix()),
		Signature:     GenerateTransactionSignature(c.Tripay.MerhcantCode, orderID, amount, c.Tripay.PrivateApiKey),
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", TripayBaseURL+"/transaction/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.ReadTransactionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Tripay.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.ReadTransactionResponse{}, err
	}
	defer resp.Body.Close()

	var raw model.Result

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return model.ReadTransactionResponse{}, fmt.Errorf("failed to decode raw result: %w", err)
	}

	if !raw.Success {
		return model.ReadTransactionResponse{}, fmt.Errorf("tripay responded with success=false")
	}

	var data model.ReadTransactionResponse
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return model.ReadTransactionResponse{}, fmt.Errorf("failed to unmarshal data field: %w", err)
	}

	fmt.Println("Checkout URL:", data.CheckoutUrl)
	fmt.Println("QR URL:", data.QrUrl)
	fmt.Println("Status:", data.Status)

	return data, nil
}

func (c *TripayClient) GetPaymentChannels() ([]model.ReadPaymentChannelResponse, error) {

	req, err := http.NewRequest("GET", TripayBaseURL+"/merchant/payment-channel", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	fmt.Println("Api_kewy:", c.Tripay.ApiKey)
	req.Header.Set("Authorization", "Bearer "+c.Tripay.ApiKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var raw model.Result

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode raw result: %w", err)
	}

	if !raw.Success {
		return nil, fmt.Errorf("tripay responded with success=false")
	}

	var data []model.ReadPaymentChannelResponse
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data field: %w", err)
	}
	return data, nil

}

func TicketsToItemsTr(tickets []*entity.Ticket) []model.OrderItem {
	var items []model.OrderItem
	for _, t := range tickets {
		name := "Tiket " + t.Type
		if t.Type == "passenger" && t.PassengerName != nil {
			name = fmt.Sprintf("Tiket Penumpang - %s", *t.PassengerName)
		}
		if t.Type == "vehicle" && t.LicensePlate != nil {
			name = fmt.Sprintf("Tiket Kendaraan - %s", *t.LicensePlate)
		}

		item := model.OrderItem{
			SKU:      t.Class.ClassName,
			Name:     name,         // e.g. "Passenger", "Vehicle"
			Price:    int(t.Price), // pastikan tipe t.Price sesuai (float32 ke int)
			Quantity: 1,            // satu tiket per entri
			// url: optional (tidak wajib), bisa ditambahkan jika ada
		}
		items = append(items, item)
	}
	return items
}
