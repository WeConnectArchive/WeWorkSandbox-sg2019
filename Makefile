protoc -I api/accounts/ api/accounts/accounts.proto --go_out=plugins=grpc:api/accounts
protoc -I api/billing/ api/billing/billing.proto --go_out=plugins=grpc:api/billing
protoc -I api/payments/ api/payments/payments.proto --go_out=plugins=grpc:api/payments