package utils

import (
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

func SendEmail(To string, subject string, body string) error {

	senderEmail := os.Getenv("SMTP_SENDER")
	Port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", "no-reply <"+senderEmail+">")
	msg.SetHeader("To", To)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	n := gomail.NewDialer(os.Getenv("SMTP_HOST"), Port, senderEmail, os.Getenv("SMTP_PASSWORD"))

	if err := n.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
