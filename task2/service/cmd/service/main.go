package main

import (
	"log"

	config "github.com/panttony/golang-course/service/config"
	http_client "github.com/panttony/golang-course/service/internal/adapter"
	grpc_server "github.com/panttony/golang-course/service/internal/controller/grpc"
	grpc_handler "github.com/panttony/golang-course/service/internal/usecase"
)

func main() {
	cfg := config.LoadServiceConfig()

	githubClient := http_client.New(cfg.GitHubAPIBaseURL, cfg.GitHubTimeout)

	handler := grpc_handler.New(githubClient)
	server := grpc_server.New(cfg.GRPCAddr, handler)

	if err := server.Run(); err != nil {
		log.Fatalf("service stopped with error: %v", err)
	}
}
