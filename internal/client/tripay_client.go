package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"eticket-api/config"
	"eticket-api/internal/common/httpclient"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	TripayBaseURL = "https://mockservicetripay.proudpebble-66ad7167.southeastasia.azurecontainerapps.io/api"
)

type TripayClient struct {
	HTTPClient *httpclient.HTTP
	Tripay     *config.Tripay
}

func NewTripayClient(httpClient *httpclient.HTTP, congig *config.Config) *TripayClient {
	return &TripayClient{
		HTTPClient: httpClient,
		Tripay:     &congig.Tripay,
	}
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

func GenerateTransactionSignature(merchantCode, merchantRef string, amount int, privateKey string) string {
	message := merchantCode + merchantRef + fmt.Sprintf("%d", amount)
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func TicketToItem(ticket *domain.Ticket) domain.OrderItem {
	name := "Tiket " + ticket.Type
	if ticket.Type == "passenger" && ticket.PassengerName != "" {
		name = fmt.Sprintf("Tiket Penumpang - %s", ticket.PassengerName)
	}
	if ticket.Type == "vehicle" && ticket.LicensePlate != nil {
		name = fmt.Sprintf("Tiket Kendaraan - %s", *ticket.LicensePlate)
	}
	return domain.OrderItem{
		SKU:      ticket.Class.ClassName,
		Name:     name,
		Price:    int(ticket.Price),
		Quantity: 1,
	}
}

func (c *TripayClient) CreatePayment(payload *domain.TransactionRequest) (*domain.Transaction, error) {
	payload.Signature = GenerateTransactionSignature(c.Tripay.MerhcantCode, payload.MerchantRef, payload.Amount, c.Tripay.PrivateApiKey)
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", TripayBaseURL+"/transaction/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Tripay.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Handle timeout or connection error
		if os.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("tripay timeout: %w", err)
		}
		return nil, fmt.Errorf("tripay connection error: %w", err)
	}
	defer resp.Body.Close()

	// Print raw response body to console
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	fmt.Println("Raw response body:", string(rawBody))

	var raw model.Result
	if err := json.Unmarshal(rawBody, &raw); err != nil {
		return nil, fmt.Errorf("failed to decode raw result: %w", err)
	}

	if !raw.Success {
		return nil, fmt.Errorf("tripay responded with success=false")
	}

	var data *domain.Transaction
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data field: %w", err)
	}

	fmt.Println("Checkout URL:", data.CheckoutUrl)
	fmt.Println("QR URL:", data.QrUrl)
	fmt.Println("Status:", data.Status)

	return data, nil
}

func (c *TripayClient) GetPaymentChannels() ([]*domain.PaymentChannel, error) {

	req, err := http.NewRequest("GET", TripayBaseURL+"/merchant/payment-channel", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	fmt.Println("Api_key:", c.Tripay.ApiKey)
	req.Header.Set("Authorization", "Bearer "+c.Tripay.ApiKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Handle timeout or connection error
		if os.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("tripay timeout: %w", err)
		}
		return nil, fmt.Errorf("tripay connection error: %w", err)
	}
	defer resp.Body.Close()

	var raw model.Result
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode raw result: %w", err)
	}

	if !raw.Success {
		return nil, fmt.Errorf("tripay responded with success=false")
	}

	var data []*domain.PaymentChannel
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data field: %w", err)
	}
	return data, nil

}

func (c *TripayClient) GetTransactionDetail(reference string) (*domain.Transaction, error) {
	req, err := http.NewRequest("GET", TripayBaseURL+"/transaction/detail", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	q := req.URL.Query()
	q.Add("reference", reference)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.Tripay.ApiKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Handle timeout or connection error
		if os.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("tripay timeout: %w", err)
		}
		return nil, fmt.Errorf("tripay connection error: %w", err)
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
	var detail *domain.Transaction
	if err := json.Unmarshal(raw.Data, &detail); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction detail: %w", err)
	}

	return detail, nil
}

// func (c *TripayClient) HandleCallback(r *http.Request) error {
// 	rawBody, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read body: %w", err)
// 	}
// 	r.Body = io.NopCloser(bytes.NewBuffer(rawBody)) // reset body if needed again

// 	signature := r.Header.Get("X-Callback-Signature")
// 	if VerifyCallbackSignature(c.Tripay.PrivateApiKey, rawBody, signature) {
// 		return fmt.Errorf("invalid callback signature")
// 	}

// 	var payload model.WriteCallbackRequest
// 	if err := json.Unmarshal(rawBody, &payload); err != nil {
// 		return fmt.Errorf("invalid JSON payload: %w", err)
// 	}

// 	// Lanjutkan dengan penyimpanan / update status pembayaran, dsb
// 	fmt.Println("Callback payload verified:", payload)
// 	return nil
// }
