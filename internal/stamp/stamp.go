// Package stamp is respnsible for stamping the last run time to BigQuery.
package stamp

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/puppetlabs/cat-team-github-metrics/internal/bigqueryclient"
	"github.com/puppetlabs/cat-team-github-metrics/internal/configuration"
	"github.com/puppetlabs/cat-team-github-metrics/internal/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bq bigqueryclient.BigQueryClient

func Run() error {
	var err error

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	err = configuration.InitConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	config := configuration.Config

	log.Info().Msg("Initializing BigQuery client")
	ctx := context.Background()
	bq, err = bigqueryclient.NewBigQueryClient(ctx, config.BigQueryProjectID, config.BigQueryDatasetName)
	if err != nil {
		return fmt.Errorf("failed to create BigQuery client: %w", err)
	}

	lastRunTime := time.Now()
	log.Info().Msgf("Last run time: %s", lastRunTime.Format(time.RFC3339))

	log.Info().Msg("Uploading last run time.")
	err = bq.Insert(config.LastRunTable, metrics.LastRunMetric{
		LastRunTime:    lastRunTime,
		CollectionTime: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to insert issue metrics: %w", err)
	}

	return nil
}
