package accounts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/weworksandbox/sg2019/api/billing"
	"github.com/weworksandbox/sg2019/api/payments"
)

// Server represents the gRPC server
type Server struct {
	BillingClient  billing.BillingClient
	PaymentsClient payments.PaymentsClient
}

// PayInvoice returns invoice and error
func (s *Server) PayInvoice(ctx context.Context, req *Invoice) (*Invoice, error) {
	invoiceID := req.Id

	// process payment for invoice
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	r1, err := s.PaymentsClient.MakePayment(timeoutCtx, &payments.PaymentRequest{})
	if err != nil {
		return nil, fmt.Errorf("payment error: %v", err)
	}
	if !r1.Paid {
		return nil, errors.New("failed to pay")
	}
	log.Printf("Paid: %t", r1.Paid)

	// mark invoice as paid
	timeoutCtx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()
	r2, err := s.BillingClient.MarkInvoicePaid(timeoutCtx, &billing.Invoice{Id: invoiceID})
	if err != nil {
		return nil, fmt.Errorf("mark invoice paid error: %v", err)
	}
	if !r2.Paid {
		return nil, errors.New("failed to mark invoice as paid")
	}
	log.Printf("Paid: %t", r2.Paid)

	req.Paid = r2.Paid

	return req, nil
}
