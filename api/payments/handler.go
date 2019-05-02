package payments

import (
	context "context"
)

// Server represents the gRPC server
type Server struct {
}

// MakePayment returns invoice and error
func (s *Server) MakePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	return &PaymentResponse{Paid: true}, nil
}
