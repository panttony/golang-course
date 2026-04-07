package grpcclient

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"

	http_v1 "github.com/panttony/golang-course/api/gen/http/gateway"
	proto_v1 "github.com/panttony/golang-course/api/gen/proto/service"
)

type Client struct {
	cli proto_v1.GitHubServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	return &Client{
		cli: proto_v1.NewGitHubServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (http_v1.RepositoryResponse, error) {
	req := proto_v1.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	}

	resp, err := c.cli.GetRepository(ctx, &req)
	if err != nil {
		return http_v1.RepositoryResponse{}, err
	}

	return repositoryFromProto(resp), nil
}

func repositoryFromProto(resp *proto_v1.GetRepositoryResponse) http_v1.RepositoryResponse {
	return http_v1.RepositoryResponse{
		Owner:           resp.GetOwner(),
		Name:            resp.GetName(),
		Description:     resp.GetDescription(),
		StargazersCount: resp.GetStargazersCount(),
		ForksCount:      resp.GetForksCount(),
		CreatedAt:       resp.GetCreatedAt().AsTime(),
	}
}
