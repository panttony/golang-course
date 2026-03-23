package http

import (
	"context"
	"fmt"

	gen "github.com/example/github-two-services/internal/gateway/transport/http/gen"
)

type service interface {
	GetRepository(ctx context.Context, owner, repo string) (gen.RepositoryResponse, error)
}

type Handler struct {
	repoService service
}

func NewHandler(repoService service) *Handler {
	return &Handler{
		repoService: repoService,
	}
}

func (h *Handler) GetRepository(
	ctx context.Context,
	req gen.GetRepositoryRequestObject,
) (gen.GetRepositoryResponseObject, error) {
	resp, err := h.repoService.GetRepository(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from grpc service: %w", err)
	}

	return gen.GetRepository200JSONResponse(resp), nil
}
