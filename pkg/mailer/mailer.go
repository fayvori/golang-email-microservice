package pkg

import (
	"crypto/tls"
	"go-email/config"

	gomail "gopkg.in/gomail.v2"
)

// email dialer
func NewMailDialer(cfg *config.Config) *gomail.Dialer {
	d := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Password)

	//nolint
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}
