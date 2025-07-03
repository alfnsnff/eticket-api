package job

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
)

type Email struct {
	To      string
	Subject string
	Body    string
}

type EmailJob struct {
	Queue  chan Email
	Mailer mailer.Mailer
	Log    logger.Logger
}

// Pastikan EmailJob implementasi mailer.Mailer
var _ mailer.Mailer = (*EmailJob)(nil)

func NewEmailJobQueue(mailer mailer.Mailer, log logger.Logger) *EmailJob {
	q := &EmailJob{
		Queue:  make(chan Email, 100), // buffer sesuai kebutuhan
		Mailer: mailer,
		Log:    log,
	}
	go q.worker()
	return q
}

func (q *EmailJob) worker() {
	q.Log.Info("Worker started") // tambahkan ini
	for job := range q.Queue {
		q.Log.Info("Sending email", "to", job.To, "subject", job.Subject)
		q.Mailer.SendAsync(job.To, job.Subject, job.Body)
	}
}

// Fungsi untuk push job (async)
func (q *EmailJob) SendAsync(to, subject, body string) {
	q.Queue <- Email{To: to, Subject: subject, Body: body}
}

func (q *EmailJob) Send(to, subject, body string) error {
	q.Queue <- Email{To: to, Subject: subject, Body: body}
	return nil
}
