//Package metrics contains methods that are responsible for mapping responses to metrics
//that can be sent to BigQuery.
package metrics

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/chelnak/puppet-github-metrics/internal/githubclient"
)

// ReleaseMetric is a struct that implements the ValueSaver interface for saving to BigQuery
type ReleaseMetric struct {
	Repository     string
	Name           string
	PublishedAt    time.Time
	CollectionTime time.Time
}

func (i *ReleaseMetric) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Repository":     i.Repository,
		"Name":           i.Name,
		"PublishedAt":    i.PublishedAt.Unix(),
		"CollectionTime": i.CollectionTime.Unix(),
	}, "", nil
}

func GetReleaseMetrics(org string, repo string) ([]ReleaseMetric, error) {
	client, err := githubclient.NewGitHubClient()
	if err != nil {
		return nil, err
	}

	metrics, err := client.GetLatestRelease(context.Background(), org, repo)
	if err != nil {
		return nil, err
	}

	return mapReleaseMetrics(metrics), nil
}

func mapReleaseMetrics(metrics []githubclient.Release) []ReleaseMetric {
	var mapped []ReleaseMetric
	for _, m := range metrics {
		mapped = append(mapped, ReleaseMetric{
			Repository:     m.Repository,
			Name:           m.Name,
			PublishedAt:    m.PublishedAt,
			CollectionTime: time.Now(),
		})
	}
	return mapped
}
