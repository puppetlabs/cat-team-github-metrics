package metrics

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/puppetlabs/cat-team-github-metrics/internal/githubclient"
)

// PullRequestMetric is a struct that implements the ValueSaver interface for saving to BigQuery
type PullRequestMetric struct {
	Repository     string
	Number         int
	Title          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Labels         []string
	Author         string
	State          string
	Merged         bool
	CollectionTime time.Time
}

func (i *PullRequestMetric) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Repository":     i.Repository,
		"Number":         i.Number,
		"Title":          i.Title,
		"CreatedAt":      i.CreatedAt.Unix(),
		"UpdatedAt":      i.UpdatedAt.Unix(),
		"Labels":         i.Labels,
		"Author":         i.Author,
		"State":          i.State,
		"Merged":         i.Merged,
		"CollectionTime": i.CollectionTime.Unix(),
	}, "", nil
}

func GetPullRequestMetrics(org string, repo string) ([]PullRequestMetric, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, err
	}

	metrics, err := client.GetPullRequests(context.Background(), org, repo)
	if err != nil {
		return nil, err
	}

	return mapPullRequestMetrics(metrics), nil
}

func mapPullRequestMetrics(metrics []githubclient.PullRequest) []PullRequestMetric {
	var mapped []PullRequestMetric
	for _, m := range metrics {
		mapped = append(mapped, PullRequestMetric{
			Repository:     m.Repository,
			Number:         m.Number,
			Title:          m.Title,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
			Labels:         m.Labels,
			Author:         m.Author,
			State:          m.State,
			Merged:         m.Merged,
			CollectionTime: time.Now(),
		})
	}
	return mapped
}
