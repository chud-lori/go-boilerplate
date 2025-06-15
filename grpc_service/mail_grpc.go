package grpc_service

import (
	"context"
	"log"
	"time"

	pb "github.com/chud-lori/go-boilerplate/proto"

	"google.golang.org/grpc"
)

func SendGrpcMail(email string, message string) {
	_ = SendGrpcMailWithOpts(email, message, grpc.WithInsecure(), grpc.WithBlock())

	// conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }

	// defer conn.Close()

	// c := pb.NewMailClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// // try with empty email to trigger error
	// r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
	// if err != nil {
	// 	// convert err to grpc status
	// 	if grpcStatus, ok := status.FromError(err); ok {
	// 		log.Printf("Error code: %v", grpcStatus.Code())
	// 		log.Printf("Error message: %v", grpcStatus.Message())

	// 		// get additional error details from metadata
	// 		for key, value := range grpcStatus.Proto().GetDetails() {
	// 			log.Printf("Error detail - %v: %v", key, value)
	// 		}
	// 	} else {
	// 		log.Printf("Unexpected error: %v", err)
	// 	}
	// } else {
	// 	log.Printf("GRPC Success: %s", r.GetMessage())
	// }
}

func SendGrpcMailWithOpts(email string, message string, dialOpts ...grpc.DialOption) error {
	conn, err := grpc.DialContext(context.Background(), "bufnet", dialOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewMailClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
	if err != nil {
		return err
	}
	log.Printf("GRPC Success: %s", r.GetMessage())
	return nil
}
