package mailer

import (
	"eticket-api/config"
	"fmt"
	"mime"
	"net/smtp"
	"strings"
)

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTP(cfg *config.Config) *SMTP {
	return &SMTP{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}
}

func (m *SMTP) SendAsync(toEmail, subject, body string) {
	go func() {
		if err := m.Send(toEmail, subject, body); err != nil {
			// Log error here if you have a logger
			fmt.Printf("failed to send email: %v\n", err)
		}
	}()
}

// send is the actual synchronous email sender (private)
func (m *SMTP) Send(toEmail, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	encodedSubject := EncodeRFC2047(subject)
	encodedFrom := EncodeRFC2047(m.From)

	msg := strings.Join([]string{
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		fmt.Sprintf("From: %s", encodedFrom),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Subject: %s", encodedSubject),
		"", // Blank line between headers and body
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	return smtp.SendMail(addr, auth, m.From, []string{toEmail}, []byte(msg))
}

// encodeRFC2047 safely encodes a string for use in email headers (Subject, From, etc.)
func EncodeRFC2047(s string) string {
	if IsASCII(s) {
		return s
	}
	return mime.QEncoding.Encode("UTF-8", s)
}

func IsASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}
