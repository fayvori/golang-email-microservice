package test

import (
	"go-email/internal/mailer"
	"go-email/internal/models"
	mail "go-email/pkg/mailer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailSending_EmailSendingNoError(t *testing.T) {
	d := mail.NewMailDialer(cfg)
	mailer := mailer.NewMailer(d)

	email := &models.Email{
		From:        cfg.Smtp.User,
		To:          []string{"alexemailtestingtwo@yahoo.com"},
		ContentType: "text/plain",
		Subject:     "Testing",
		Body:        "Test email",
	}

	err := mailer.SendEmails(email)
	require.Nil(t, err)
	require.NoError(t, err)
}
