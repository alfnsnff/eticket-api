package mailer

import (
	"eticket-api/config"
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPMailer(cfg *config.Configuration) *SMTPMailer {
	return &SMTPMailer{
		Host:     cfg.SMTPMailer.Host,
		Port:     cfg.SMTPMailer.Port,
		Username: cfg.SMTPMailer.Username,
		Password: cfg.SMTPMailer.Password,
		From:     cfg.SMTPMailer.From,
	}
}
func (m *SMTPMailer) Send(toEmail, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	msg := "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += fmt.Sprintf("From: %s\r\n", m.From)
	msg += fmt.Sprintf("To: %s\r\n", toEmail)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "\r\n" + body

	// Auth method
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	// Send email
	return smtp.SendMail(addr, auth, m.From, []string{toEmail}, []byte(msg))
}
