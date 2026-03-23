package githubapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	domain "github.com/example/github-two-services/internal/service/domain"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type RepoResponse struct {
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

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetByFullName(ctx context.Context, owner, repo string) (domain.RepoResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.RepoResponse{}, fmt.Errorf("create github request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.RepoResponse{}, fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.RepoResponse{}, fmt.Errorf(messageError(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.RepoResponse{}, fmt.Errorf("read github response: %w", err)
	}

	var dto domain.RepoResponse
	if err := json.Unmarshal(body, &dto); err != nil {
		return domain.RepoResponse{}, fmt.Errorf("decode github response: %w", err)
	}

	return dto, nil
}

func messageError(code int) string {
	switch code {
	case http.StatusMovedPermanently:
		return http.StatusText(code)
	case http.StatusForbidden:
		return http.StatusText(code)
	case http.StatusNotFound:
		return http.StatusText(code)
	default:
	}

	return "unknown error"
}
