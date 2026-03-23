package query

import (
	"context"
	"net/http"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	github_v1 "github.com/pantonny/golang-course/internal/gen/proto/github/v1"
)

type githubReader interface {
	GetByFullName(ctx context.Context, owner, repo string) (dto *github_v1.GetRepositoryResponse, statusCode int, err error)
}

type githubRequestHandler struct {
	github_v1.UnimplementedGitHubServiceServer
	reader githubReader
}

func NewGetRepositoryQuery(reader githubReader) *githubRequestHandler {
	return &githubRequestHandler{
		reader: reader,
	}
}

func (q *githubRequestHandler) GetRepository(ctx context.Context, req *github_v1.GetRepositoryRequest) (*github_v1.GetRepositoryResponse, error) {
	resp, statusCode, err := q.reader.GetByFullName(ctx, req.GetOwner(), req.GetRepo())

	if err != nil {
		return nil, codeToError(statusCode)
	}

	return resp, nil
}

func codeToError(code int) error {
	switch code {
	case http.StatusMovedPermanently:
		return status.Errorf(codes.NotFound, "Resource moved permanently")
	case http.StatusForbidden:
		return status.Errorf(codes.PermissionDenied, "Permission denied")
	case http.StatusNotFound:
		return status.Errorf(codes.NotFound, "Resource not found")
	default:
		return status.Errorf(codes.Internal, "Internal server error")
	}
}
