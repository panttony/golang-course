package main

import (
	"log"

	query "github.com/pantonny/golang-course/internal/service/application/query"
	github_api "github.com/pantonny/golang-course/internal/service/infrastructure/github_api"
	grpc_server "github.com/pantonny/golang-course/internal/service/transport/grpc"
	config "github.com/pantonny/golang-course/internal/shared/config"
)

func main() {
	cfg := config.LoadServiceConfig()

	githubClient := github_api.NewClient(cfg.GitHubAPIBaseURL, cfg.GitHubTimeout)
	useCase := query.NewGetRepositoryQuery(githubClient)

	server := grpc_server.New(cfg.GRPCAddr, useCase)
	if err := server.Run(); err != nil {
		log.Fatalf("service stopped with error: %v", err)
	}
}
