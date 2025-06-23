package mailer

type Mailer interface {
	Send(toEmail, subject, body string) error
}
