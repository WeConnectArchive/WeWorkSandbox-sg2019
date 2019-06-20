it-prepare:
	@ echo Prepare env. It is blocking operation...
	@ go run services/accounts/main.go

it-run:
	@ echo Run tests...
	@ go test -count=1 -v ./test/integration/...

generate:
	protoc -I api/accounts/ api/accounts/accounts.proto --go_out=plugins=grpc:api/accounts
	protoc -I api/billing/ api/billing/billing.proto --go_out=plugins=grpc:api/billing
	protoc -I api/payments/ api/payments/payments.proto --go_out=plugins=grpc:api/payments