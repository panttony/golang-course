PROTO_DIR := api/proto
PROTO_OUT := api/gen/proto

OPENAPI_SPEC := api/openapi/gateway.yaml
OPENAPI_CONFIG := api/openapi/config/oapi-codegen.yaml
SWAGGER_DOC_OUT := gateway/internal/http/docs/gateway.yaml

OPENAPI_TOOL := $(shell go env GOPATH)/bin/oapi-codegen
PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc

.PHONY: generate generate-openapi generate-proto install-tools up down build clean-gen

generate: install-tools generate-openapi generate-proto

generate-docker: generate-openapi generate-proto

install-tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate-openapi:
	mkdir -p $(dir $(OPENAPI_SPEC))
	mkdir -p $(dir $(SWAGGER_DOC_OUT))
	$(OPENAPI_TOOL) -config $(OPENAPI_CONFIG) $(OPENAPI_SPEC)
	cp $(OPENAPI_SPEC) $(SWAGGER_DOC_OUT)

generate-proto:
	mkdir -p $(PROTO_OUT)
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc \
		-I $(PROTO_DIR) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/service/service.proto

up:
	docker compose up --build

down:
	docker compose down -v

build:
	go build -o bin/gateway ./gateway/cmd/...
	go build -o bin/service ./service/cmd/...

clean-gen:
	rm -rf api/gen
	rm -f gateway/internal/http/gen/*.go
	rm -f $(SWAGGER_DOC_OUT)
