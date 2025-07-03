package mailer

type Mailer interface {
	SendAsync(to, subject, body string)
	Send(to, subject, body string) error
}
