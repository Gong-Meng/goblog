package mail

import (
	"goblog/pkg/logger"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) bool {
	m := gomail.NewMessage()
	m.SetHeader("From", "local@localhost")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("172.17.0.4", 1025, "", "")

	if err := d.DialAndSend(m); err != nil {
		logger.LogError(err)
		return false
	}

	return true
}
