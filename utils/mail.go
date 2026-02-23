package utils

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
	username = "dhaaranisakthii@gmail.com"
	password = "tiybscsimienxtch" // Gmail App Password
)

func SendCredentials(to, name, user, pass, role string) error {
	log.Println("📨 SMTP → Sending mail to:", to)

	msg := fmt.Sprintf(
		"Subject: PySellers Login Credentials\r\n"+
			"\r\n"+
			"Hi %s,\n\n"+
			"Your account has been created.\n\n"+
			"Login: http://localhost:9090\n"+
			"Username: %s\n"+
			"Password: %s\n"+
			"Role: %s\n\n"+
			"Please change password after login.\n\n"+
			"HR Team",
		name, user, pass, role,
	)

	auth := smtp.PlainAuth("", username, password, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		username,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Println("❌ SMTP ERROR:", err)
		return err
	}

	log.Println("📬 SMTP → Mail accepted by server")
	return nil
}
