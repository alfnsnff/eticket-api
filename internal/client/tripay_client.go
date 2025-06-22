package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"eticket-api/config"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"
	"io"
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

func (c *TripayClient) CreatePayment(method string, amount int, name string, email string, phone string, orderID string, items []*domain.Ticket) (model.ReadTransactionResponse, error) {
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ReadTransactionResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}
	fmt.Println("Raw response:", string(bodyBytes))

	// Reset resp.Body for further decoding
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

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
	fmt.Println("Api_key:", c.Tripay.ApiKey)
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

func (c *TripayClient) GetTransactionDetail(reference string) (*model.ReadTransactionResponse, error) {
	// Build request URL with query parameter
	req, err := http.NewRequest("GET", TripayBaseURL+"/transaction/detail", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add reference as query param
	q := req.URL.Query()
	q.Add("reference", reference)
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.Tripay.ApiKey)

	// Send request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse base Tripay response
	var raw model.Result
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !raw.Success {
		return nil, fmt.Errorf("tripay responded with success=false")
	}

	// Unmarshal `raw.Data` into expected detail type
	var detail model.ReadTransactionResponse
	if err := json.Unmarshal(raw.Data, &detail); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction detail: %w", err)
	}

	return &detail, nil
}

func (c *TripayClient) HandleCallback(r *http.Request) error {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(rawBody)) // reset body if needed again

	signature := r.Header.Get("X-Callback-Signature")
	if VerifyCallbackSignature(c.Tripay.PrivateApiKey, rawBody, signature) {
		return fmt.Errorf("invalid callback signature")
	}

	var payload model.WriteCallbackRequest
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}

	// Lanjutkan dengan penyimpanan / update status pembayaran, dsb
	fmt.Println("Callback payload verified:", payload)
	return nil
}

func VerifyCallbackSignature(private_api_key string, raw_body []byte, signature string) bool {
	expected := GenerateCallbackSignature(private_api_key, raw_body)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func GenerateCallbackSignature(privateKey string, rawBody []byte) string {
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write(rawBody)
	return hex.EncodeToString(h.Sum(nil))
}

func TicketsToItemsTr(tickets []*domain.Ticket) []model.OrderItem {
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
