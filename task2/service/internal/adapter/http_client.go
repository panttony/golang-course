package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	github_v1 "github.com/panttony/golang-course/api/gen/proto/service"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type repoResponse struct {
	Owner       Owner     `json:"owner"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stargazers  int64     `json:"stargazers_count"`
	Forks       int64     `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type Owner struct {
	Login string `json:"login"`
}

func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetByFullName(ctx context.Context, owner, repo string) (*github_v1.GetRepositoryResponse, int, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, -1, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	dto := &repoResponse{}
	if err := json.Unmarshal(body, dto); err != nil {
		return nil, resp.StatusCode, err
	}

	return &github_v1.GetRepositoryResponse{
		Owner:           dto.Owner.Login,
		Name:            dto.Name,
		Description:     dto.Description,
		StargazersCount: dto.Stargazers,
		ForksCount:      dto.Forks,
		CreatedAt:       timestamppb.New(dto.CreatedAt),
	}, resp.StatusCode, nil
}
