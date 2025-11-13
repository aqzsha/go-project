package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	From     string
	Password string
	Host     string
	Port     string
}

func (m *Mailer) SendMail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", m.From, to, subject, body)

	auth := smtp.PlainAuth("", m.From, m.Password, m.Host)
	return smtp.SendMail(addr, auth, m.From, []string{to}, []byte(msg))
}
