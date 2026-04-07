package usecase

import (
	"context"

	http_v1 "github.com/panttony/golang-course/api/gen/http/gateway"
)

type service interface {
	GetRepository(ctx context.Context, owner, repo string) (http_v1.RepositoryResponse, error)
}

type Handler struct {
	repoService service
}

func New(repoService service) *Handler {
	return &Handler{
		repoService: repoService,
	}
}

func (h *Handler) GetRepository(
	ctx context.Context,
	req http_v1.GetRepositoryRequestObject,
) (http_v1.GetRepositoryResponseObject, error) {
	resp, err := h.repoService.GetRepository(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, err
	}

	return http_v1.GetRepository200JSONResponse(resp), nil
}
