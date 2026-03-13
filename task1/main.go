package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	urlSep        = "/"
	url           = "https://api.github.com/repos"
	commonTimeout = 5 * time.Second
)

var usage = `
Usage: qwe <argument>
	argument:	github url or owner/repo
Example:
	qwe https://github.com/panttony/golang-course
	or
	qwe panttony/golang-course
`

type GitHubRepo struct {
	Owner       Owner  `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type Owner struct {
	Login string `json:"login"`
}

func (g GitHubRepo) String() string {
	return fmt.Sprintf("Owner: %s\n"+
		"Name of repository: %s\n"+
		"Description: %s\n"+
		"Count of stars: %d\n"+
		"Count of forks: %d\n"+
		"Created at: %s\n",
		g.Owner.Login, g.Name, g.Description, g.Stargazers, g.Forks, g.CreatedAt,
	)	
}

func cli(args []string) error {
	args = args[1:]
	arg := args[0]

	splitArgs := strings.Split(arg, urlSep)
	if len(splitArgs) < 2 {
		return fmt.Errorf("incorrect url format: %s", arg)
	}

	owner, repo := splitArgs[len(splitArgs)-2], splitArgs[len(splitArgs)-1]
	if owner == "" || repo == "" {
		return fmt.Errorf("incorrect arguments (owner=%s, repo=%s)", owner, repo)
	}

	req, err := buildReq(owner, repo)
	if err != nil {
		return err
	}

	c := newClient(commonTimeout)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(messageError(resp.StatusCode))
	}

	gh, err := readResponse(resp)
	if err != nil {
		return err
	}

	fmt.Print(gh)

	return nil
}

func readResponse(resp *http.Response) (*GitHubRepo, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := &GitHubRepo{}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}

	return v, nil
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

func newClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func buildReq(owner, repo string) (resp *http.Request, err error) {
	url := fmt.Sprintf("%s/%s/%s", url, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request is not created: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	return req, nil
}

func main() {
	if len(os.Args) < 2 {
		panic(usage)
	}

	if err := cli(os.Args); err != nil {
		fmt.Printf("%v\n", err)

		os.Exit(1)
	}
}
