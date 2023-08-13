package pkg

import (
	"crypto/tls"
	config "go-email/config"

	gomail "gopkg.in/gomail.v2"
)

// email dialer
func NewMailDialer(cfg *config.Config) *gomail.Dialer {
	d := gomail.NewDialer(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.User, cfg.Smtp.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}
