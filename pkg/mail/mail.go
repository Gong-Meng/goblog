package mail

import (
	"goblog/pkg/logger"

	"gopkg.in/gomail.v2"
)

func SendMail() bool {
	m := gomail.NewMessage()
	m.SetHeader("From", "local@localhost")
	m.SetHeader("To", "mengmeng@123.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewDialer("172.17.0.4", 1025, "", "")

	if err := d.DialAndSend(m); err != nil {
		logger.LogError(err)
		return false
	}

	return true
}
