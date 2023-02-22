package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddr   = "smtp.gmail.com"
	smtpAuthServer = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error
}

type GmailSender struct {
	name              string
	fromEmailAddr     string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddr, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddr:     fromEmailAddr,
		fromEmailPassword: fromEmailPassword,
	}
}

func (g *GmailSender) SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", g.name, g.fromEmailAddr)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.Subject = subject
	e.HTML = []byte(content)

	for _, attachFile := range attachFiles {
		if _, err := e.AttachFile(attachFile); err != nil {
			return fmt.Errorf("failed to attach file: %w", err)
		}
	}
	auth := smtp.PlainAuth("", g.fromEmailAddr, g.fromEmailPassword, smtpAuthAddr)
	return e.Send(smtpAuthServer, auth)
}
