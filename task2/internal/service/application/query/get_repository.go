package query

import (
	"context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	
	github_v1 "github.com/example/github-two-services/internal/gen/proto/github/v1"
	domain "github.com/example/github-two-services/internal/service/domain"
)

var errInternal = status.Errorf(codes.Internal, "Internal server error")

type githubReader interface {
	GetByFullName(ctx context.Context, owner, repo string) (domain.RepoResponse, error)
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
	httpResp, err := q.reader.GetByFullName(ctx, req.GetOwner(), req.GetRepo())

	if err != nil {
		return nil, errInternal
	}

	return httpToProto(httpResp), nil
}

func httpToProto(resp domain.RepoResponse) *github_v1.GetRepositoryResponse {
	return &github_v1.GetRepositoryResponse{
		Owner:           resp.Owner.Login,
		Name:            resp.Name,
		Description:     resp.Description,
		StargazersCount: resp.Stargazers,
		ForksCount:      resp.Forks,
		CreatedAt:       timestamppb.New(resp.CreatedAt),
	}
}
