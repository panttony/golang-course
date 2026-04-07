package main

import (
	"log"

	config "github.com/panttony/golang-course/gateway/config"
	grpc_client "github.com/panttony/golang-course/gateway/internal/adapter"
	http_server "github.com/panttony/golang-course/gateway/internal/http"
	repo_info "github.com/panttony/golang-course/gateway/internal/usecase"
)

func main() {
	cfg := config.LoadGatewayConfig()

	client, err := grpc_client.New(cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("cannot create grpc client: %v", err)
	}

	handler := repo_info.New(client)
	server := http_server.New(cfg.HTTPAddr, handler)

	if err := server.Run(); err != nil {
		log.Fatalf("gateway stopped with error: %v", err)
	}
}
