package domain

import "time"

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
