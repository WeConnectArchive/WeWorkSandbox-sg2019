package billing

import (
	context "context"
	fmt "fmt"
	"log"
	"net"

	"github.com/weworksandbox/sg2019/pkg/mock"
	grpc "google.golang.org/grpc"
)

// NewMockServer creates a new mock billing gRPC server
func NewMockServer(port int, messageChan chan interface{}) {
	// create a listener on TCP port 50052
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50052))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := MockServer{messageChan: messageChan}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	RegisterBillingServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// MockServer represents the mocked gRPC server
type MockServer struct {
	messageChan chan interface{}
}

// MarkInvoicePaid returns invoice and error
func (m *MockServer) MarkInvoicePaid(ctx context.Context, req *Invoice) (*Invoice, error) {
	response := mock.GetInterface(m.messageChan)
	if response == nil {
		log.Fatal("Test case did not program a mock response")
	}
	v, ok := response.(*Invoice)
	if !ok {
		return nil, fmt.Errorf("Mock response is the wrong type: %+v", response)
	}
	m.messageChan <- req
	return v, nil
}
