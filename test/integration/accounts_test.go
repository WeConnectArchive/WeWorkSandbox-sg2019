package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	api "github.com/weworksandbox/sg2019/api/accounts"
	"github.com/weworksandbox/sg2019/api/billing"
	"github.com/weworksandbox/sg2019/api/payments"
	"github.com/weworksandbox/sg2019/pkg/mock"
	"google.golang.org/grpc"
)

const (
	accountsAddress = "0.0.0.0:50051"
	paymentsPort    = 50052
	billingPort     = 50053
)

var (
	accountsClient  api.AccountsClient
	billingChannel  chan interface{}
	paymentsChannel chan interface{}
)

func TestMain(m *testing.M) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(accountsAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	accountsClient = api.NewAccountsClient(conn)

	// Run downstream mocks
	billingChannel = make(chan interface{}, 100)
	go billing.NewMockServer(billingPort, billingChannel)

	paymentsChannel = make(chan interface{}, 100)
	go payments.NewMockServer(paymentsPort, paymentsChannel)

	// Run tests
	os.Exit(m.Run())
}

func TestPayment(t *testing.T) {

	// Configure Mocks (Set expeted outputs from downstream services)
	billingChannel <- &billing.Invoice{Id: 1, Paid: true}
	paymentsChannel <- &payments.PaymentResponse{Paid: true}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := accountsClient.PayInvoice(ctx, &api.Invoice{Id: 1})
	if err != nil {
		log.Fatalf("Payment error: %v", err)
	}
	if !r.Paid {
		log.Fatal("Failed to pay")
	}
	log.Printf("Paid: %t", r.Paid)

	// Validate Accounts Server made a request to the Payments Service
	request := mock.GetInterface(paymentsChannel)
	if request == nil {
		log.Fatal("Accounts service did not make a request to Payments Service")
	}
	pr, ok := request.(*payments.PaymentRequest)
	if !ok {
		log.Fatalf("Payments Service recieved the wrong type of request: %v", pr)
	}

	// Validate Accounts Server made a request to the Billing Service
	request = mock.GetInterface(billingChannel)
	if request == nil {
		log.Fatal("Accounts Service did not make a request to Billing Service")
	}
	br, ok := request.(*billing.Invoice)
	if !ok {
		log.Fatalf("Billing Service recieved the wrong type of request: %v", br)
	}
	// Validate that Accounts Server tried to pay the correct invoice
	if br.Id != 1 {
		log.Fatal("Accounts Service tried to mark the wrong invoice as paid")
	}
}

func TestPaymentFailure(t *testing.T) {

	// Configure Mocks (Set expeted output for downstream services)
	paymentsChannel <- fmt.Errorf("This is a negative test")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := accountsClient.PayInvoice(ctx, &api.Invoice{Id: 1})
	if err != nil {
		log.Fatalf("Payment error: %v", err)
	}
	// We now expect the payment to have failed
	if r.Paid {
		log.Fatal("Successfully paid but expected failure")
	}
	log.Printf("Failed to pay as expected")

	// Validate Accounts Server made a request to the Payments Service
	request := mock.GetInterface(paymentsChannel)
	if request == nil {
		log.Fatal("Accounts service did not make a request to Payments Service")
	}
	pr, ok := request.(*payments.PaymentRequest)
	if !ok {
		log.Fatalf("Payments Service recieved the wrong type of request: %v", pr)
	}

	// Validate Accounts Server did not make a request to the Billing Service
	request = mock.GetInterface(billingChannel)
	if request != nil {
		log.Fatalf("Accounts Service made a request to Billing Service: %+v", request)
	}
}
