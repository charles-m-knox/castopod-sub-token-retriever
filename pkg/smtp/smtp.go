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
	// If true, emails won't be sent. Good for testing.
	Test bool `json:"testing"`
}

// SendEmail sends an email to the desired recipient.
func (s *SMTPConfig) SendEmail(to []string, subject, body string, from string) error {
	if len(to) < 1 {
		return fmt.Errorf("please provide at least 1 'to' address")
	}

	auth := smtp.PlainAuth("", from, s.Password, s.Server)

	msg := fmt.Sprintf("To: %v\r\nSubject: %v\r\n\r\n%v", to[0], subject, body)

	if s.Test {
		fmt.Printf("To (all): %v\n\n Msg: %v", to, msg)
		return nil
	}

	err := smtp.SendMail(s.Addr, auth, from, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
