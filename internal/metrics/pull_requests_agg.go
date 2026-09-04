package metrics

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/puppetlabs/cat-team-github-metrics/internal/githubclient"
)

type PullRequestAggregatedMetric struct {
	Repository     string
	Count          int
	CollectionTime time.Time
}

func (i *PullRequestAggregatedMetric) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Repository":     i.Repository,
		"Count":          i.Count,
		"CollectionTime": i.CollectionTime.Unix(),
	}, "", nil
}

func GetPullRequestAggregatedMetrics(org string, repo string) (PullRequestAggregatedMetric, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return PullRequestAggregatedMetric{}, err
	}
	return getPullRequestAggregatedMetrics(client, org, repo)
}

func getPullRequestAggregatedMetrics(client githubclient.GitHubClient, org, repo string) (PullRequestAggregatedMetric, error) {
	metrics, err := client.GetOpenPullRequestCount(context.Background(), org, repo)
	if err != nil {
		return PullRequestAggregatedMetric{}, err
	}
	return PullRequestAggregatedMetric{
		Repository:     metrics.Repository,
		Count:          metrics.Count,
		CollectionTime: time.Now(),
	}, nil
}
