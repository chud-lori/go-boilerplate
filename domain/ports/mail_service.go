package ports

import "context"

type MailService interface {
	SendSignInNotification(ctx context.Context, email, text string) error
}
