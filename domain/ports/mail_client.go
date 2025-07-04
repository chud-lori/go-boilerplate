package ports

import "context"

type MailClient interface {
	SendMail(ctx context.Context, email string, message string) error
}
