package ports

type MailClient interface {
	SendMail(email string, message string) error
}
