package metrics //nolint:dupl

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/puppetlabs/cat-team-github-metrics/internal/githubclient"
)

type IssueAggregatedMetric struct {
	Repository     string
	Count          int
	CollectionTime time.Time
}

func (i *IssueAggregatedMetric) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Repository":     i.Repository,
		"Count":          i.Count,
		"CollectionTime": i.CollectionTime.Unix(),
	}, "", nil
}

func GetIssueAggregatedMetrics(org string, repo string) (IssueAggregatedMetric, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return IssueAggregatedMetric{}, err
	}
	return getIssueAggregatedMetrics(client, org, repo)
}

func getIssueAggregatedMetrics(client githubclient.GitHubClient, org, repo string) (IssueAggregatedMetric, error) {
	metrics, err := client.GetOpenIssueCount(context.Background(), org, repo)
	if err != nil {
		return IssueAggregatedMetric{}, err
	}
	return IssueAggregatedMetric{
		Repository:     metrics.Repository,
		Count:          metrics.Count,
		CollectionTime: time.Now(),
	}, nil
}
