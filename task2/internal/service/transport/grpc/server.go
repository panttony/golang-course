package grpc

import (
	"context"
	"fmt"
	"net"

	grpc "google.golang.org/grpc"

	github_v1 "github.com/example/github-two-services/internal/gen/proto/github/v1"
)

type githubApiHandler interface {
	GetRepository(ctx context.Context, req *github_v1.GetRepositoryRequest) (*github_v1.GetRepositoryResponse, error)
}

type Server struct {
	addr    string
	srv     *grpc.Server
	handler githubApiHandler
}

func New(addr string, handler githubApiHandler) *Server {
	srv := grpc.NewServer()
	srv.RegisterService(&github_v1.GitHubService_ServiceDesc, handler)

	return &Server{addr: addr, srv: srv, handler: handler}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen grpc: %w", err)
	}
	defer lis.Close()

	return s.srv.Serve(lis)
}
