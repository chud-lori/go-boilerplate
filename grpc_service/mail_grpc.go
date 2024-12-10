package grpc_service

import (
	"context"
	"log"
	"time"

	pb "github.com/chud-lori/go-boilerplate/mail_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func SendGrpcMail(email string, message string) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewMailClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// try with empty email to trigger error
	r, err := c.SendMail(ctx, &pb.MailRequest{Email: email, Message: message})
	if err != nil {
		// convert err to grpc status
		if grpcStatus, ok := status.FromError(err); ok {
			log.Printf("Error code: %v", grpcStatus.Code())
			log.Printf("Error message: %v", grpcStatus.Message())

			// get additional error details from metadata
			for key, value := range grpcStatus.Proto().GetDetails() {
				log.Printf("Error detail - %s: %s", key, value)
			}
		} else {
			log.Printf("Unexpected error: %v", err)
		}
	} else {
		log.Printf("GRPC Success: %s", r.GetMessage())
	}
}
