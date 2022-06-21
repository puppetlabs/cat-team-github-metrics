//Package metrics contains methods that are responsible for mapping responses to metrics
//that can be sent to BigQuery.
package metrics

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/chelnak/cat-team-github-metrics/internal/githubclient"
)

// IssueMetric is a struct that implements the ValueSaver interface for saving to BigQuery
type IssueMetric struct {
	Repository     string
	Number         int
	Title          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Labels         []string
	Author         string
	State          string
	CollectionTime time.Time
}

func (i *IssueMetric) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Repository":     i.Repository,
		"Number":         i.Number,
		"Title":          i.Title,
		"CreatedAt":      i.CreatedAt.Unix(),
		"UpdatedAt":      i.UpdatedAt.Unix(),
		"Labels":         i.Labels,
		"Author":         i.Author,
		"State":          i.State,
		"CollectionTime": i.CollectionTime.Unix(),
	}, "", nil
}

func GetIssueMetrics(org string, repo string) ([]IssueMetric, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, err
	}

	metrics, err := client.GetIssues(context.Background(), org, repo)
	if err != nil {
		return nil, err
	}

	return mapIssueMetrics(metrics), nil
}

func mapIssueMetrics(metrics []githubclient.Issue) []IssueMetric {
	var mapped []IssueMetric
	for _, m := range metrics {
		mapped = append(mapped, IssueMetric{
			Repository:     m.Repository,
			Number:         m.Number,
			Title:          m.Title,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
			Labels:         m.Labels,
			Author:         m.Author,
			State:          m.State,
			CollectionTime: time.Now(),
		})
	}
	return mapped
}
