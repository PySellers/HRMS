package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendCredentials sends login details to employee email
func SendCredentials(to, name, user, pass, role string) error {
	// ==============================
	// LOAD SMTP FROM ENV
	// ==============================
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Validate required env variables
	if smtpHost == "" || smtpPort == "" || username == "" || password == "" {
		log.Println("❌ SMTP environment variables not set properly")
		log.Println("📧 Email would have been sent to:", to)
		log.Println("📧 Credentials - Username:", user, "Password:", pass)
		// Don't return error, just log it so the student creation still works
		return nil
	}

	log.Println("📨 SMTP → Sending mail to:", to)

	// ==============================
	// EMAIL CONTENT
	// ==============================
	subject := "Subject: PySellers Login Credentials\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	// HTML email body
	body := fmt.Sprintf(`
        <html>
        <body style="font-family: Arial, sans-serif; padding: 20px;">
            <h2 style="color: #3498db;">Welcome to PySellers!</h2>
            <p>Hi <strong>%s</strong>,</p>
            <p>Your account has been created as a <strong>%s</strong>.</p>
           
            <div style="background: #f4f6f9; padding: 20px; border-radius: 5px; margin: 20px 0;">
                <p style="margin: 5px 0;"><strong>Login URL:</strong> <a href="http://localhost:9090">http://localhost:9090</a></p>
                <p style="margin: 5px 0;"><strong>Username:</strong> %s</p>
                <p style="margin: 5px 0;"><strong>Password:</strong> %s</p>
            </div>
           
            <p><strong>Important:</strong> Please change your password after first login.</p>
            <p>If you have any questions, please contact the HR team.</p>
           
            <hr style="border: 1px solid #eee; margin: 20px 0;">
            <p style="color: #666; font-size: 12px;">This is an automated message, please do not reply.</p>
        </body>
        </html>
    `, name, role, user, pass)

	msg := []byte(subject + mime + body)

	// ==============================
	// SMTP AUTH
	// ==============================
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		username,
		[]string{to},
		msg,
	)

	if err != nil {
		log.Println("❌ SMTP ERROR:", err)
		log.Println("📧 Email would have been sent to:", to)
		log.Println("📧 Credentials - Username:", user, "Password:", pass)
		return err
	}

	log.Println("📬 SMTP → Mail accepted by server")
	return nil
}
