package cstr

import (
	"fmt"
	"net/smtp"
	"strings"

	uuid "git.cmcode.dev/cmcode/uuid"
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

// NewUUID returns a double uuid with the hyphens removed, leading to a
// 64-character string. Castopod appears to use something like this for its
// token.
func NewUUID() string {
	return strings.ReplaceAll(fmt.Sprintf("%v%v", uuid.New(), uuid.New()), "-", "")
}
