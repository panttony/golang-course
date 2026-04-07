package usecase

import (
	"context"
	"net/http"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	github_v1 "github.com/panttony/golang-course/api/gen/proto/service"
)

type githubCollector interface {
	GetByFullName(ctx context.Context, owner, repo string) (dto *github_v1.GetRepositoryResponse, statusCode int, err error)
}

type githubRequestHandler struct {
	github_v1.UnimplementedGitHubServiceServer
	collector githubCollector
}

func New(collector githubCollector) *githubRequestHandler {
	return &githubRequestHandler{
		collector: collector,
	}
}

func (q *githubRequestHandler) GetRepository(ctx context.Context, req *github_v1.GetRepositoryRequest) (*github_v1.GetRepositoryResponse, error) {
	resp, statusCode, err := q.collector.GetByFullName(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		return nil, codeToError(statusCode)
	}

	return resp, nil
}

func codeToError(code int) error {
	switch code {
	case http.StatusMovedPermanently, http.StatusNotFound:
		return status.Error(codes.NotFound, "repository not found")
	case http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, "authentication required")
	case http.StatusForbidden:
		return status.Error(codes.PermissionDenied, "access denied by github api")
	case http.StatusTooManyRequests:
		return status.Error(codes.ResourceExhausted, "github api rate limit exceeded")
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return status.Error(codes.Unavailable, "github api is unavailable")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
