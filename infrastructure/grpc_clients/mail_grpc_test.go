package grpc_clients_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/infrastructure/grpc_clients"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	pb "github.com/chud-lori/go-boilerplate/proto"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func freshBreaker() *gobreaker.CircuitBreaker[[]byte] {
	var st gobreaker.Settings
	st.Name = "TestGrpcMailClient"
	st.MaxRequests = 3
	st.Interval = 60 * time.Second
	st.Timeout = 10 * time.Second
	return gobreaker.NewCircuitBreaker[[]byte](st)
}

// Test helper for creating client with custom breaker
func newGrpcMailClientWithBreaker(conn grpc.ClientConnInterface, breaker *gobreaker.CircuitBreaker[[]byte]) *grpc_clients.GrpcMailClient {
	return &grpc_clients.GrpcMailClient{
		Conn:    conn,
		Breaker: breaker,
	}
}

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

// Custom mock server for circuit breaker tests
// Always fails
type alwaysFailMailServer struct {
	pb.UnimplementedMailServer
}

func (s *alwaysFailMailServer) SendMail(ctx context.Context, req *pb.MailRequest) (*pb.MailReply, error) {
	return nil, status.Errorf(codes.Internal, "simulated failure")
}

// Fails N times, then succeeds
type failThenSucceedMailServer struct {
	pb.UnimplementedMailServer
	failCount int
	failLimit int
}

func (s *failThenSucceedMailServer) SendMail(ctx context.Context, req *pb.MailRequest) (*pb.MailReply, error) {
	if s.failCount < s.failLimit {
		s.failCount++
		return nil, status.Errorf(codes.Internal, "simulated failure")
	}
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
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	lis, s, conn := setupTestServer(t)
	defer conn.Close()
	defer s.Stop()    // Explicitly stop the server
	defer lis.Close() // Close the listener

	client := newGrpcMailClientWithBreaker(conn, freshBreaker())

	email := "test@example.com"
	message := "Hello, this is a test email."

	err := client.SendMail(ctx, email, message)
	if err != nil {
		t.Errorf("SendMail failed unexpectedly: %v", err)
	}
}

func TestGrpcMailClient_SendMail_InvalidInput(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	lis, s, conn := setupTestServer(t)
	defer conn.Close()
	defer s.Stop()    // Explicitly stop the server
	defer lis.Close() // Close the listener

	client := newGrpcMailClientWithBreaker(conn, freshBreaker())

	// Test case: empty email
	err := client.SendMail(ctx, "", "Some message")
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
	err = client.SendMail(ctx, "test@example.com", "")
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

func TestGrpcMailClient_CircuitBreaker_OpenState(t *testing.T) {
	baseCtx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterMailServer(s, &alwaysFailMailServer{})

	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Printf("gRPC server exited with error: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(baseCtx, 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet for test: %v", err)
	}
	defer conn.Close()
	defer s.Stop()
	defer lis.Close()

	client := newGrpcMailClientWithBreaker(conn, freshBreaker())

	// Make multiple requests to trigger circuit breaker
	for i := 0; i < 5; i++ {
		err := client.SendMail(baseCtx, "test@example.com", "Test message")
		if err == nil {
			t.Errorf("Expected error on request %d", i)
		}
	}

	// After multiple failures, circuit breaker should be open
	err = client.SendMail(baseCtx, "test@example.com", "Test message")
	if err == nil {
		t.Error("Expected circuit breaker to be open and return error")
	}
}

func TestGrpcMailClient_CircuitBreaker_Recovery(t *testing.T) {
	baseCtx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterMailServer(s, &failThenSucceedMailServer{failLimit: 2})

	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Printf("gRPC server exited with error: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(baseCtx, 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet for test: %v", err)
	}
	defer conn.Close()
	defer s.Stop()
	defer lis.Close()

	client := newGrpcMailClientWithBreaker(conn, freshBreaker())

	// Make a few failing requests
	for i := 0; i < 2; i++ {
		err := client.SendMail(baseCtx, "test@example.com", "Test message")
		if err == nil {
			t.Errorf("Expected error on request %d", i)
		}
	}

	// Make a successful request
	err = client.SendMail(baseCtx, "test@example.com", "Test message")
	if err != nil {
		t.Errorf("Expected successful request, got error: %v", err)
	}
}

func TestGrpcMailClient_ConcurrentRequests(t *testing.T) {
	baseCtx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	lis, s, conn := setupTestServer(t)
	defer conn.Close()
	defer s.Stop()
	defer lis.Close()

	client := newGrpcMailClientWithBreaker(conn, freshBreaker())

	// Make concurrent requests
	const numRequests = 5
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			err := client.SendMail(baseCtx, "test@example.com", "Test message")
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		err := <-results
		if err != nil {
			t.Errorf("Concurrent request %d failed: %v", i, err)
		}
	}
}
