package qr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type InvoiceItems struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type InvoiceRequest struct {
	ExternalID  string `json:"external_id"`
	PayerEmail  string `json:"payer_email"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	// Items              []InvoiceItems `json:"items"`
	SuccessRedirectURL string `json:"success_redirect_url"`
	FailureRedirectURL string `json:"failure_redirect_url"`
}

type InvoiceResponse struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
}

type QRISRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int    `json:"amount"`
	Type        string `json:"type"` // "DYNAMIC"
	CallbackURL string `json:"callback_url"`
	Description string `json:"description"`
}

type QRISResponse struct {
	ID       string `json:"id"`
	QRString string `json:"qr_string"`
	QRURL    string `json:"qr_url"`
	Status   string `json:"status"`
	Amount   int    `json:"amount"`
}

type TicketToInvoiceItem struct {
	ClassName string
	Type      string
	Price     float32
}

// func MapTicketsToInvoiceItems(tickets []*entity.Ticket) []InvoiceItems {
// 	var ticketDTOs []TicketToInvoiceItem
// 	err := copier.Copy(&ticketDTOs, &tickets)
// 	if err != nil {
// 		return nil
// 	}

// 	var invoiceItems []InvoiceItems
// 	for _, t := range ticketDTOs {
// 		item := InvoiceItems{
// 			Class: ClassItem{
// 				ClassName: t.ClassName,
// 				Type:      t.ClassCode,
// 			},
// 			Price: t.Price,
// 		}
// 		invoiceItems = append(invoiceItems, item)
// 	}

// 	return invoiceItems
// }

func CreateInvoice(payload InvoiceRequest) (InvoiceResponse, error) {
	url := "https://api.xendit.co/v2/invoices"
	apiKey := "xnd_development_5jmjfvncIPR6hLvvXbfoJoINBkChlXRVWCyQaIl7sMR0DlSeZRpFfQB8VyP6AX" // Ganti dengan sandbox API Key Anda

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

	var result InvoiceResponse
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("Invoice URL:", result.InvoiceURL)
	return result, nil
}

func CreateQRIS(payload QRISRequest) (QRISResponse, error) {
	url := "https://api.xendit.co/qr_codes"
	apiKey := "xnd_development_5jmjfvncIPR6hLvvXbfoJoINBkChlXRVWCyQaIl7sMR0DlSeZRpFfQB8VyP6AX" // Ganti dengan sandbox API Key Anda

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

	var result QRISResponse
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("QRIS URL:", result.QRURL)
	fmt.Println("QR STRING:", result.QRString)
	return result, nil
}
