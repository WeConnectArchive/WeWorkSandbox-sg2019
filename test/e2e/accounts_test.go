package test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	api "github.com/weworksandbox/sg2019/api/accounts"
	"google.golang.org/grpc"
)

const (
	accountsAddress = "0.0.0.0:50051"
)

var (
	accountsClient api.AccountsClient
)

func TestMain(m *testing.M) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(accountsAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	accountsClient = api.NewAccountsClient(conn)

	// Run tests
	os.Exit(m.Run())
}

func TestPayment(t *testing.T) {
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
}
