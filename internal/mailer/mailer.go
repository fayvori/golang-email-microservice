package mailer

import (
	"go-email/internal/models"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(dialer *gomail.Dialer) *Mailer {
	return &Mailer{dialer: dialer}
}

func (m *Mailer) SendEmails(email *models.Email) error {
	emailDialer := gomail.NewMessage()
	emailDialer.SetHeader("From", email.From)
	emailDialer.SetHeader("To", email.To...)
	emailDialer.SetHeader("Subject", email.Subject)

	emailDialer.SetBody(email.ContentType, email.Body)

	return m.dialer.DialAndSend(emailDialer)
}
