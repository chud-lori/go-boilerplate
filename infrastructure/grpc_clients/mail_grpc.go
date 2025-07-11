package grpc_clients

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	pb "github.com/chud-lori/go-boilerplate/proto"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"github.com/sony/gobreaker"
)

var grpcMailBreaker *gobreaker.CircuitBreaker

func init() {
	grpcMailBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "GrpcMailClient",
		MaxRequests: 3,
		Interval:    60 * 1e9, // 60s
		Timeout:     10 * 1e9, // 10s
	})
}

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

	_, err := grpcMailBreaker.Execute(func() (interface{}, error) {
		r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
		if err != nil {
			logger.WithError(err).Error("could not send mail")
			return nil, err
		}
		logger.Infof("GRPC Success: %s", r.GetMessage())
		return nil, nil
	})
	return err
}
