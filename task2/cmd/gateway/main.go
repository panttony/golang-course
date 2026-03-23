package main

import (
	"log"

	grpc_client "github.com/pantonny/golang-course/internal/gateway/adapter/grpc_client"
	http_transport "github.com/pantonny/golang-course/internal/gateway/transport/http"
	config "github.com/pantonny/golang-course/internal/shared/config"
)

func main() {
	cfg := config.LoadGatewayConfig()

	client, err := grpc_client.New(cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("cannot create grpc client: %v", err)
	}

	strictHandler := http_transport.NewHandler(client)
	server := http_transport.NewServer(cfg.HTTPAddr, strictHandler)

	if err := server.Run(); err != nil {
		log.Fatalf("gateway stopped with error: %v", err)
	}
}
