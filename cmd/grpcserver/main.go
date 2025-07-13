package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"time"

	"github.com/chud-lori/go-boilerplate/config"
	"github.com/chud-lori/go-boilerplate/internal/utils"
	pb "github.com/chud-lori/go-boilerplate/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mailServer struct {
	pb.UnimplementedMailServer
}

func (s *mailServer) SendMail(ctx context.Context, req *pb.MailRequest) (*pb.MailReply, error) {
	if req.GetEmail() == "" || req.GetMessage() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email or message cannot be empty")
	}

	auth := smtp.PlainAuth("", config.Mail.User, config.Mail.Pass, config.Mail.Host)

	smtpAddr := fmt.Sprintf("%s:%d", config.Mail.Host, config.Mail.Port)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Testing\r\n\r\n%s", req.Email, req.Message))

	err := smtp.SendMail(smtpAddr, auth, config.Mail.From, []string{req.Email}, msg)

	if err != nil {
		log.Printf("SMTP send error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send email: %v", err)
	}

	log.Printf("✅ Sent mail to %s", req.Email)
	return &pb.MailReply{
		Message: fmt.Sprintf(`Success send "%s" to %s`, req.Message, req.Email),
	}, nil
}

func main() {
	config.LoadMailConfig()

	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("❌ Failed to listen on %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMailServer(grpcServer, &mailServer{})

	log.Println("📨 gRPC Mail server running on port", port)

	// Start the gRPC server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("❌ gRPC server failed to serve: %v", err)
		}
	}()

	// ========== Graceful Shutdown ==========
	wait := utils.GracefullShutdown(context.Background(), 30*time.Second, map[string]utils.Operation{
		"grpc-server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			log.Println("gRPC server GracefulStop completed.")
			return nil
		},
	})

	<-wait
	log.Println("🚀 gRPC Server exited.")
}

