package billing

import (
	context "context"
)

// Server represents the gRPC server
type Server struct {
}

// MarkInvoicePaid returns invoice and error
func (s *Server) MarkInvoicePaid(ctx context.Context, req *Invoice) (*Invoice, error) {
	req.Paid = true
	return req, nil
}
