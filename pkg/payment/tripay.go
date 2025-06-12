package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"eticket-api/internal/entity"
	"fmt"
	"net/http"
	"time"
)

type OrderItem struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type Data struct {
	Reference   string `json:"reference"`
	CheckoutUrl string `json:"checkout_url"`
	QrUrl       string `json:"qr_url"`
	Status      string `json:"status"`
}

type TripayRequest struct {
	Method        string      `json:"method"`
	MerchantRef   string      `json:"merchant_ref"`
	Amount        int         `json:"amount"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	OrderItems    []OrderItem `json:"order_items"`
	CallbackUrl   string      `json:"callback_url"`
	ReturnUrl     string      `json:"return_url"`
	ExpiredTime   int         `json:"expired_time"`
	Signature     string      `json:"signature"`
}

type TripayResponse struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

func GenerateTransactionSignature(merchantCode, merchantRef string, amount int, privateKey string) string {
	message := merchantCode + merchantRef + fmt.Sprintf("%d", amount)
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	TripayApiKey       = "DEV-guDYuHQycRa3KbrIarTncvUXPddhZ9hAEX8uebvE"
	TripayMerchantCode = "T41200"
	TripayPrivateKey   = "3B4fT-2lG29-itTeU-dtKFM-twvHS"
	TripayBaseURL      = "https://tripay.co.id/api-sandbox"
)

func CreateTripayPayment(method string, amount int, name string, email string, orderID string, items []*entity.Ticket) (TripayResponse, error) {

	payload := TripayRequest{
		Method:        method,
		MerchantRef:   orderID,
		Amount:        amount,
		CustomerName:  name,
		CustomerEmail: email,
		OrderItems:    TicketsToItemsTr(items),
		CallbackUrl:   "https://example.com/callback",
		ReturnUrl:     "https://example.com/success",
		ExpiredTime:   int(time.Now().Add(30 * time.Minute).Unix()),

		Signature: GenerateTransactionSignature(TripayMerchantCode, orderID, amount, TripayPrivateKey),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", TripayBaseURL+"/transaction/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+TripayApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TripayResponse{}, err
	}
	defer resp.Body.Close()

	var result TripayResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if !result.Success {
		return TripayResponse{}, fmt.Errorf("gagal membuat transaksi")
	}

	fmt.Println("Checkout URL:", result.Data.CheckoutUrl)
	fmt.Println("QR URL:", result.Data.QrUrl)
	fmt.Println("Status:", result.Data.Status)

	return result, err
}

func TicketsToItemsTr(tickets []*entity.Ticket) []OrderItem {
	var items []OrderItem
	for _, t := range tickets {
		name := "Tiket " + t.Type
		if t.Type == "passenger" && t.PassengerName != nil {
			name = fmt.Sprintf("Tiket Penumpang - %s", *t.PassengerName)
		}
		if t.Type == "vehicle" && t.LicensePlate != nil {
			name = fmt.Sprintf("Tiket Kendaraan - %s", *t.LicensePlate)
		}

		item := OrderItem{
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
