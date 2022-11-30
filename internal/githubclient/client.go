// Package githubclient is responsible for interacting with the GitHub API
package githubclient

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	GetIssues(ctx context.Context, owner string, repo string) ([]Issue, error)
	GetOpenIssueCount(ctx context.Context, owner string, repo string) (IssueCount, error)
	GetPullRequests(ctx context.Context, owner string, repo string) ([]PullRequest, error)
	GetOpenPullRequestCount(ctx context.Context, owner string, repo string) (PullRequestCount, error)
	GetLatestRelease(ctx context.Context, owner string, repo string) ([]Release, error)
}

type githubClient struct {
	v4 *githubv4.Client
}

func NewGitHubClient() (GitHubClient, error) {
	httpClient, err := newOauthHTTPClient()
	if err != nil {
		return nil, err
	}

	v4Client := githubv4.NewClient(httpClient)

	client := &githubClient{
		v4: v4Client,
	}

	return client, nil
}

func newOauthHTTPClient() (*http.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("GITHUB_TOKEN is not set")
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	return httpClient, nil
}
