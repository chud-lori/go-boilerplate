package services

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

type MailServiceImpl struct {
	ports.MailClient
}

func (s *MailServiceImpl) SendSignInNotification(ctx context.Context, email, text string) error {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	err := s.MailClient.SendMail(email, text)
	if err != nil {
		logger.WithError(err).Error("Error send notif")
		return err
	}

	return nil
}
