package mailer

import (
	"fmt"
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
	gm := gomail.NewMessage()
	gm.SetHeader("From", email.From)
	gm.SetHeader("To", email.To...)
	gm.SetHeader("Subject", email.Subject)

	gm.SetBody(email.ContentType, email.Body)

	fmt.Println(email)

	return m.dialer.DialAndSend(gm)
}
