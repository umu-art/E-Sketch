package api

type MailApi interface {
	SendConfirmationEmail(email, confirmationLink string, token string) error
}
