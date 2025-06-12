package payment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"eticket-api/internal/entity"
	"fmt"
	"net/http"
	"strings"
)

type Items struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type XenditRequest struct {
	ExternalID         string  `json:"external_id"`
	PayerEmail         string  `json:"payer_email"`
	Description        string  `json:"description"`
	Amount             int     `json:"amount"`
	Items              []Items `json:"items"`
	SuccessRedirectURL string  `json:"success_redirect_url"`
	FailureRedirectURL string  `json:"failure_redirect_url"`
}

type XenditResponse struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
}

// type TicketToItem struct {
// 	ClassName string
// 	Type      string
// 	Price     float32
// }

func TicketsToItems(tickets []*entity.Ticket) []Items {
	var items []Items
	for _, t := range tickets {
		name := "Tiket " + t.Type
		if t.Type == "passenger" && t.PassengerName != nil {
			name = fmt.Sprintf("Tiket Penumpang - %s", *t.PassengerName)
		}
		if t.Type == "vehicle" && t.LicensePlate != nil {
			name = fmt.Sprintf("Tiket Kendaraan - %s", *t.LicensePlate)
		}

		item := Items{
			Name:     name,
			Category: strings.Title(t.Type), // e.g. "Passenger", "Vehicle"
			Quantity: 1,                     // satu tiket per entri
			Price:    int(t.Price),          // pastikan tipe t.Price sesuai (float32 ke int)
			// url: optional (tidak wajib), bisa ditambahkan jika ada
		}
		items = append(items, item)
	}
	return items
}

func CreateXenditPayment(amount int, name string, email string, orderID string, items []*entity.Ticket) (XenditResponse, error) {
	url := "https://api.xendit.co/v2/invoices"
	apiKey := "xnd_development_5jmjfvncIPR6hLvvXbfoJoINBkChlXRVWCyQaIl7sMR0DlSeZRpFfQB8VyP6AX" // Ganti dengan sandbox API Key Anda

	payload := XenditRequest{
		ExternalID:         orderID,
		PayerEmail:         email,
		Description:        "Pembayaran tiket kapal untuk Booking #" + orderID,
		Amount:             amount,
		Items:              TicketsToItems(items),
		SuccessRedirectURL: "https://yourdomain.com/payment-success",
		FailureRedirectURL: "https://yourdomain.com/payment-failed",
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result XenditResponse
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("Invoice URL:", result.InvoiceURL)
	return result, nil
}
