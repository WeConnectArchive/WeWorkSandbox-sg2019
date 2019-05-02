package main

import (
	"fmt"
	"log"
	"net"

	api "github.com/weworksandbox/sg2019/api/accounts"
	"github.com/weworksandbox/sg2019/api/billing"
	"github.com/weworksandbox/sg2019/api/payments"
	"google.golang.org/grpc"
)

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := api.Server{}

	// setup billing client
	conn1, err := grpc.Dial("0.0.0.0:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()
	s.BillingClient = billing.NewBillingClient(conn1)

	// setup payments client
	conn2, err := grpc.Dial("0.0.0.0:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn2.Close()
	s.PaymentsClient = payments.NewPaymentsClient(conn2)

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	api.RegisterAccountsServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
