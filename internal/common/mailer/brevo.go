package mailer

import (
	"bytes"
	"encoding/json"
	"eticket-api/config"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Brevo struct {
	APIKey   string
	From     string
	FromName string
}

func NewBrevo(cfg *config.Config) *Brevo {
	return &Brevo{
		APIKey:   cfg.Brevo.APIKey,
		From:     cfg.Brevo.From,
		FromName: cfg.Brevo.Name,
	}
}

func (b *Brevo) SendAsync(to, subject, body string) {
	go func() {
		err := b.Send(to, subject, body)
		if err != nil {
			fmt.Printf("[Brevo] failed to send email to %s: %v\n", to, err)
		}
	}()
}

func (b *Brevo) Send(to, subject, body string) error {
	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  b.FromName,
			"email": b.From,
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject":     subject,
		"htmlContent": body,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", b.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("email send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("email send failed, status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// Ensure Brevo implements Mailer interface
var _ Mailer = (*Brevo)(nil)
