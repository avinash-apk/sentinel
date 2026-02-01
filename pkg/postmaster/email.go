package postmaster

import (
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	From     string
	Password string
	Host     string
	Port     string
}

func NewEmailSender(user, password string) *EmailSender {
	return &EmailSender{
		From:     user,
		Password: password,
		Host:     "smtp.gmail.com",
		Port:     "587",
	}
}

func (e *EmailSender) Send(destination string, body string) error {
	// Standard SMTP Auth
	auth := smtp.PlainAuth("", e.From, e.Password, e.Host)
	
	// Create the email headers and body
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Sentinel Notification\r\n"+
		"\r\n"+
		"%s\r\n", destination, body))

	addr := e.Host + ":" + e.Port
	return smtp.SendMail(addr, auth, e.From, []string{destination}, msg)
}