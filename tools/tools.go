package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway" //nolint
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2" //nolint
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc" //nolint
	_ "google.golang.org/protobuf/cmd/protoc-gen-go" //nolint
)
