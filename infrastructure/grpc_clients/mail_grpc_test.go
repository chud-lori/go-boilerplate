package grpc_clients_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/infrastructure/grpc_clients"
	pb "github.com/chud-lori/go-boilerplate/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// mockMailServer is a mock implementation of your MailService gRPC server
type mockMailServer struct {
	pb.UnimplementedMailServer // Embed this to ensure forward compatibility
}

// SendMail implements the SendMail method of the MailService
func (s *mockMailServer) SendMail(ctx context.Context, req *pb.MailRequest) (*pb.MailReply, error) {
	// Simulate some server-side logic
	if req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email cannot be empty")
	}
	if req.GetMessage() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message cannot be empty")
	}

	// Simulate success
	return &pb.MailReply{Message: "Email sent successfully to " + req.GetEmail()}, nil
}

// setupTestServer creates a new bufconn.Listener, starts a gRPC server on it,
// and returns the listener, the server, and a client connection.
func setupTestServer(t *testing.T) (*bufconn.Listener, *grpc.Server, *grpc.ClientConn) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterMailServer(s, &mockMailServer{})

	// Start the server in a goroutine
	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			// Log unexpected errors, but ignore ErrServerStopped which happens on graceful shutdown
			log.Printf("gRPC server exited with error: %v", err)
		}
	}()

	// Establish client connection to the bufconn listener
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Short timeout for connection attempt
	defer cancel()                                                          // Ensure the context is cancelled

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial() // This tells gRPC to use our in-memory buffer listener
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Use insecure creds for bufconn
		grpc.WithBlock(), // This ensures the dial completes before proceeding, safe with context timeout
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet for test: %v", err)
	}

	return lis, s, conn
}

func TestGrpcMailClient_SendMail_Success(t *testing.T) {
	lis, s, conn := setupTestServer(t)
	defer conn.Close()
	defer s.Stop()    // Explicitly stop the server
	defer lis.Close() // Close the listener

	client := grpc_clients.NewGrpcMailClient(conn)

	email := "test@example.com"
	message := "Hello, this is a test email."

	err := client.SendMail(email, message)
	if err != nil {
		t.Errorf("SendMail failed unexpectedly: %v", err)
	}
}

func TestGrpcMailClient_SendMail_InvalidInput(t *testing.T) {
	lis, s, conn := setupTestServer(t)
	defer conn.Close()
	defer s.Stop()    // Explicitly stop the server
	defer lis.Close() // Close the listener

	client := grpc_clients.NewGrpcMailClient(conn)

	// Test case: empty email
	err := client.SendMail("", "Some message")
	if err == nil {
		t.Error("SendMail did not return an error for empty email")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Errorf("Error is not a gRPC status error: %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error, got %v", st.Code())
	}

	// Test case: empty message
	err = client.SendMail("test@example.com", "")
	if err == nil {
		t.Error("SendMail did not return an error for empty message")
	}
	st, ok = status.FromError(err)
	if !ok {
		t.Errorf("Error is not a gRPC status error: %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error, got %v", st.Code())
	}
}
