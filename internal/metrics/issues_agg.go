package metrics

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

	metrics, err := client.GetOpenIssueCount(context.Background(), org, repo)
	if err != nil {
		return IssueAggregatedMetric{}, err
	}

	issueMetric := IssueAggregatedMetric{
		Repository:     metrics.Repository,
		Count:          metrics.Count,
		CollectionTime: time.Now(),
	}

	return issueMetric, nil
}
