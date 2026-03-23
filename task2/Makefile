APP_NAME := github-two-services
PROTO_DIR := api/proto
PROTO_OUT := internal/gen/proto
OPENAPI_SPEC := api/openapi/gateway.yaml
OPENAPI_CONFIG := api/openapi/oapi-codegen.yaml
OPENAPI_TOOL := $(shell go env GOPATH)/bin/oapi-codegen
PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc

.PHONY: help generate generate-openapi generate-proto generate-tools fmt lint test run-gateway run-service up down build clean-gen

help:
	@echo "make generate       - install generators and generate both HTTP and gRPC code"
	@echo "make generate-tools - install oapi-codegen and protoc Go plugins"
	@echo "make generate-openapi - generate REST server/client code from OpenAPI"
	@echo "make generate-proto - generate gRPC code via protoc"
	@echo "make up             - start docker compose"
	@echo "make down           - stop docker compose"
	@echo "make build          - build both services"
	@echo "make clean-gen      - remove generated code"

generate: generate-tools generate-openapi generate-proto

generate-tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate-openapi:
	mkdir -p $(dir $(OPENAPI_CONFIG))
	$(OPENAPI_TOOL) -config $(OPENAPI_CONFIG) $(OPENAPI_SPEC)

generate-proto:
	mkdir -p $(PROTO_OUT)
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc \
		-I $(PROTO_DIR) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/github/v1/github.proto

up:
	docker compose up --build

down:
	docker compose down -v

build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/service ./cmd/service

clean-gen:
	rm -rf internal/gen
	rm -f internal/gateway/transport/http/gen/*.go
