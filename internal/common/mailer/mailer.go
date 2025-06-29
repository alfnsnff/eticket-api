package mailer

type Mailer interface {
	SendAsync(to, subject, body string)
}
