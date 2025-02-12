package api

import (
	"est-proxy/src/config"
	"fmt"
	"net/smtp"
)

type MailApi struct {
	smtpServer string
	port       string
	username   string
	password   string
}

func NewMailApi() *MailApi {
	return &MailApi{
		smtpServer: config.SMTP_SERVER,
		port:       config.SMTP_PORT,
		username:   config.SMTP_EMAIL,
		password:   config.SMTP_PASSWORD,
	}
}

func (m MailApi) SendConfirmationEmail(email, confirmationLink string) error {
	from := m.username
	subject := "Подтверждение аккаунта"

	body := fmt.Sprintf(
		"Пожалуйста, подтвердите ваш аккаунт, перейдя по следующей ссылке: %s",
		confirmationLink)

	msg := []byte("From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: " + subject + "\n\n" +
		body)

	addr := fmt.Sprintf("%s:%s", m.smtpServer, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.smtpServer)

	err := smtp.SendMail(addr, auth, from, []string{email}, msg)
	if err != nil {
		return err
	}
	return nil
}
