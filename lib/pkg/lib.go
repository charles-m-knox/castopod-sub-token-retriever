package cstr

import (
	"fmt"
	"net/smtp"
)

// SendEmail sends an email to the desired recipient.
func SendEmail(to []string, subject, body string, from string, password string, smtpServer string, smtpPort string) error {
	auth := smtp.PlainAuth("", from, password, smtpServer)

	addr := fmt.Sprintf("%v:%v", smtpServer, smtpPort)

	msg := fmt.Sprintf("To: %v\r\nSubject: %v\r\n\r\n%v", to, subject, body)

	err := smtp.SendMail(addr, auth, from, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
