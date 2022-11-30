// Package metrics contains methods that are responsible for mapping responses to metrics
// that can be sent to BigQuery.
package metrics

import (
	"time"

	"cloud.google.com/go/bigquery"
)

// LastRunMetric is a struct that implements the ValueSaver interface for saving to BigQuery
type LastRunMetric struct {
	LastRunTime    time.Time
	CollectionTime time.Time
}

func (i *LastRunMetric) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"LastRunTime":    i.LastRunTime.Unix(),
		"CollectionTime": i.CollectionTime.Unix(),
	}, "", nil
}
