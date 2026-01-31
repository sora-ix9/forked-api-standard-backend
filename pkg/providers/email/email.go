package email

import (
	"fdlp-standard-api/pkg/config"
	"strconv"

	"gopkg.in/gomail.v2"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendEmail(to []string, subject string, body string, contentType string) error
}

type emailSender struct {
	dialer *gomail.Dialer
	from   string
}

// NewEmailSender initializes and returns a new EmailSender
func NewEmailSender(cfg *config.Config) EmailSender {
	port, _ := strconv.Atoi(cfg.SMTPPort)
	dialer := gomail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUsername, cfg.SMTPPassword)

	return &emailSender{
		dialer: dialer,
		from:   cfg.SMTPFromEmail,
	}
}

// SendEmail sends an email
func (s *emailSender) SendEmail(to []string, subject string, body string, contentType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	if contentType == "" {
		contentType = "text/plain"
	}
	m.SetBody(contentType, body)

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
