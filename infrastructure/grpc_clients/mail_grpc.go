package grpc_clients

import (
	"context"
	"log"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	pb "github.com/chud-lori/go-boilerplate/proto"

	"google.golang.org/grpc"
)

type GrpcMailClient struct {
	conn grpc.ClientConnInterface
}

var _ ports.MailClient = (*GrpcMailClient)(nil)

// NewGrpcMailClient creates a new GrpcMailClient
func NewGrpcMailClient(conn grpc.ClientConnInterface) *GrpcMailClient {
	return &GrpcMailClient{conn: conn}
}

func (g *GrpcMailClient) SendMail(email string, message string) error {
	c := pb.NewMailClient(g.conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
	if err != nil {
		log.Printf("could not send mail: %v", err)
		return err
	}
	log.Printf("GRPC Success: %s", r.GetMessage())
	return nil
}

// func SendGrpcMail(email string, message string) {
// 	_ = SendGrpcMailWithOpts(email, message, grpc.WithInsecure(), grpc.WithBlock())
// }

// func SendGrpcMailWithOpts(email string, message string, dialOpts ...grpc.DialOption) error {
// 	conn, err := grpc.DialContext(context.Background(), "bufnet", dialOpts...)
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	c := pb.NewMailClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("GRPC Success: %s", r.GetMessage())
// 	return nil
// }
