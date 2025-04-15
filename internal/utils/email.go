package utils

import (
	"financierGo/config"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	cfg := config.Load()
	port, _ := strconv.Atoi(cfg.SMTP.Port)
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTP.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		cfg.SMTP.Host,
		port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
	)

	return d.DialAndSend(m)
}
