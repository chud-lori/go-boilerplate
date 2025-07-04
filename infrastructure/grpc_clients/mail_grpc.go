package grpc_clients

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	pb "github.com/chud-lori/go-boilerplate/proto"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type GrpcMailClient struct {
	conn grpc.ClientConnInterface
}

var _ ports.MailClient = (*GrpcMailClient)(nil)

func NewGrpcMailClient(conn grpc.ClientConnInterface) *GrpcMailClient {
	return &GrpcMailClient{conn: conn}
}

func (g *GrpcMailClient) SendMail(ctx context.Context, email string, message string) error {
	c := pb.NewMailClient(g.conn)

	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
	if err != nil {
		logger.WithError(err).Error("could not send mail")
		return err
	}
	logger.Infof("GRPC Success: %s", r.GetMessage())
	return nil
}
