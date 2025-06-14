package utils

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	APIKey string
	From   string
}

func NewSendGridMailer(apiKey, from string) *SendGridMailer {
	return &SendGridMailer{
		APIKey: apiKey,
		From:   from,
	}
}

func (m *SendGridMailer) Send(toEmail, subject, plainBody, htmlBody string) error {
	from := mail.NewEmail("eTicket Support", m.From)
	to := mail.NewEmail("", toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainBody, htmlBody)
	client := sendgrid.NewSendClient(m.APIKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("email failed with status: %d, body: %s", response.StatusCode, response.Body)
	}
	return nil
}
