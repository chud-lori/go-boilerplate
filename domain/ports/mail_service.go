package ports

type MailService interface {
	SendSignInNotification(email, text string) error
}
