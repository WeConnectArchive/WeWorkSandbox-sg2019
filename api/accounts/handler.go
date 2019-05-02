package accounts

import (
	context "context"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r1, err := s.PaymentsClient.MakePayment(ctx, &payments.PaymentRequest{})
	if err != nil {
		log.Fatalf("Payment error: %v", err)
	}
	if !r1.Paid {
		log.Fatal("Failed to pay")
	}
	log.Printf("Paid: %t", r1.Paid)

	// mark invoice as paid
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r2, err := s.BillingClient.MarkInvoicePaid(ctx, &billing.Invoice{Id: invoiceID})
	if err != nil {
		log.Fatalf("Mark invoice paid error: %v", err)
	}
	if !r2.Paid {
		log.Fatal("Failed to mark invoice as paid")
	}
	log.Printf("Paid: %t", r2.Paid)

	req.Paid = r2.Paid

	return req, nil
}
