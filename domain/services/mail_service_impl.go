package services

import (
	"fmt"

	"github.com/chud-lori/go-boilerplate/domain/ports"
)

type MailServiceImpl struct {
	ports.MailClient
}

func (s *MailServiceImpl) SendSignInNotification(email, text string) error {
	message := fmt.Sprintf("Notif sent")
	return s.MailClient.SendMail(email, message)
}
