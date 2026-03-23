package grpcclient

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"

	"github.com/example/github-two-services/internal/gateway/transport/http/gen"
	github_v1 "github.com/example/github-two-services/internal/gen/proto/github/v1"
)

type Client struct {
	cli github_v1.GitHubServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	return &Client{
		cli: github_v1.NewGitHubServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (gen.RepositoryResponse, error) {
	req := &github_v1.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	}

	resp, err := c.cli.GetRepository(ctx, req)
	if err != nil {
		return gen.RepositoryResponse{}, fmt.Errorf("failed to receive response from grpc server: %w", err)
	}

	return repositoryFromProto(resp), nil
}

func repositoryFromProto(resp *github_v1.GetRepositoryResponse) gen.RepositoryResponse {
	return gen.RepositoryResponse{
		Owner:           resp.GetOwner(),
		Name:            resp.GetName(),
		Description:     resp.GetDescription(),
		StargazersCount: resp.GetStargazersCount(),
		ForksCount:      resp.GetForksCount(),
		CreatedAt:       resp.GetCreatedAt().AsTime(),
	}
}
