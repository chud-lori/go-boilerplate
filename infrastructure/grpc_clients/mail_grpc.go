package grpc_clients

import (
	"context"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	pb "github.com/chud-lori/go-boilerplate/proto"
	"github.com/sirupsen/logrus"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
)

type GrpcMailClient struct {
	Conn    grpc.ClientConnInterface
	Breaker *gobreaker.CircuitBreaker[[]byte]
}

func defaultBreaker() *gobreaker.CircuitBreaker[[]byte] {
	var st gobreaker.Settings
	st.Name = "GrpcMailClient"
	st.MaxRequests = 3
	st.Interval = 60 * time.Second
	st.Timeout = 10 * time.Second
	return gobreaker.NewCircuitBreaker[[]byte](st)
}

var _ ports.MailClient = (*GrpcMailClient)(nil)

func NewGrpcMailClient(conn grpc.ClientConnInterface) *GrpcMailClient {
	return &GrpcMailClient{
		Conn:    conn,
		Breaker: defaultBreaker(),
	}
}

func (g *GrpcMailClient) SendMail(ctx context.Context, email string, message string) error {
	c := pb.NewMailClient(g.Conn)

	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	_, err := g.Breaker.Execute(func() ([]byte, error) {
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
