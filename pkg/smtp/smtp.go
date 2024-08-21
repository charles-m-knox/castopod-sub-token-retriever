package smtp

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	// SMTP server host, such as 127.0.0.1. Used for constructing the proper
	// authentication parameters.
	Server string `json:"server"`
	// SMTP server's listen address, such as 127.0.0.1:587.
	Addr string `json:"addr"`
	// SMTP server password.
	Password string `json:"password"`
}

// SendEmail sends an email to the desired recipient.
func (s *SMTPConfig) SendEmail(to []string, subject, body string, from string) error {
	auth := smtp.PlainAuth("", from, s.Password, s.Server)

	msg := fmt.Sprintf("To: %v\r\nSubject: %v\r\n\r\n%v", to, subject, body)

	err := smtp.SendMail(s.Addr, auth, from, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
