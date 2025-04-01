package impl

import (
	"est-proxy/src/config"
	"fmt"
	"github.com/toorop/go-dkim"
	"log"
	"net/smtp"
	"os"
	"time"
)

type MailApiImpl struct {
	smtpServer string
	port       string
	username   string
	password   string
	sigOptions *dkim.SigOptions
}

func NewMailApi() *MailApiImpl {
	options := dkim.NewSigOptions()

	keyBytes, err := os.ReadFile(config.SMTP_DKIM_KEY_FILE)
	if err != nil {
		log.Fatalf("Failed to read smtp cert file")
	}
	options.PrivateKey = keyBytes

	options.Domain = "e-sketch.ru"
	options.Selector = "default"
	options.SignatureExpireIn = 3600
	options.BodyLength = 0
	options.Headers = []string{"from", "to", "subject", "date", "content-type"}
	options.AddSignatureTimestamp = true
	options.Canonicalization = "relaxed/relaxed"

	return &MailApiImpl{
		smtpServer: config.SMTP_SERVER,
		port:       config.SMTP_PORT,
		username:   config.SMTP_EMAIL,
		password:   config.SMTP_PASSWORD,
		sigOptions: &options,
	}
}

func (m MailApiImpl) SendConfirmationEmail(email, confirmationLink string, token string) error {
	from := m.username
	subject := "Подтверждение аккаунта"

	body := fmt.Sprintf(
		"Пожалуйста, подтвердите ваш аккаунт, введя код %s или перейдя по ссылке %s",
		token, confirmationLink)

	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body)

	addr := fmt.Sprintf("%s:%s", m.smtpServer, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.smtpServer)

	err := dkim.Sign(&msg, *m.sigOptions)
	if err != nil {
		return fmt.Errorf("DKIM sign failed: %w", err)
	}

	err = smtp.SendMail(addr, auth, from, []string{email}, msg)
	if err != nil {
		return err
	}
	return nil
}
