package main

import (
	"log"

	grpc_client "github.com/example/github-two-services/internal/gateway/adapter/grpc_client"
	http_transport "github.com/example/github-two-services/internal/gateway/transport/http"
	config "github.com/example/github-two-services/internal/shared/config"
)

func main() {
	cfg := config.LoadGatewayConfig()

	client, err := grpc_client.New(cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("cannot create grpc client: %v", err)
	}

	strictHander := http_transport.NewHandler(client)
	server := http_transport.NewServer(cfg.HTTPAddr, strictHander)

	if err := server.Run(); err != nil {
		log.Fatalf("gateway stopped with error: %v", err)
	}
}
