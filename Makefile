# ==============================================================================
# Main commands

run:
	@go run ./cmd/email_service/main.go

build:
	@go build ./cmd/email_service/main.go

tests:
	go test -v ./test/...

clean:
	go mod tidy && go fmt ./...

lint:
	golangci-lint run \
		./config \
		./internal/... \
		./pkg/... \
		./test

proto_folder = ./pkg/proto/
openapiv2_folder = ./docs/swagger/

proto:
	protoc -I $(proto_folder) \
  	--go_out $(proto_folder)/email-service --go_opt paths=source_relative \
  	--go-grpc_out $(proto_folder)/email-service --go-grpc_opt paths=source_relative \
  	--grpc-gateway_out $(proto_folder)/email-service --grpc-gateway_opt paths=source_relative \
		--openapiv2_out $(openapiv2_folder) \
  	$(proto_folder)/mailer.proto

# ==============================================================================
# Docker
compose_up:
	@docker-compose up -f ./deployments/docker-compose.yml
