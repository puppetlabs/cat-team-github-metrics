package main

import (
	"context"
	"os"
	"time"

	"github.com/chelnak/cat-team-github-metrics/internal/bigqueryclient"
	"github.com/chelnak/cat-team-github-metrics/internal/configuration"
	"github.com/chelnak/cat-team-github-metrics/internal/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bq bigqueryclient.BigQueryClient

func main() {
	var err error

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	err = configuration.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config error")
		os.Exit(1)
	}

	config := configuration.Config

	moduleName := config.ModuleName
	moduleOwner := config.ModuleOwner

	log.Info().Msg("Initializing BigQuery client")
	ctx := context.Background()
	bq, err = bigqueryclient.NewBigQueryClient(ctx, config.BigQueryProjectID, config.BigQueryDatasetName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create BigQuery client")
		os.Exit(1)
	}

	var hasErrors bool

	log.Info().Msg("starting to collect metrics")
	log.Info().Msgf("%s: fetching issue metrics", moduleName)

	issues, err := metrics.GetIssueMetrics(moduleOwner, moduleName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch issue metrics", moduleName)
		hasErrors = true
	}

	log.Info().Msgf("%s: uploading %d issue metrics", moduleName, len(issues))
	err = bq.Insert(config.IssuesTable, issues)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to insert issue metrics", moduleName)
		hasErrors = true
	}

	log.Info().Msgf("%s: fetching pull request metrics", moduleName)
	pullRequests, err := metrics.GetPullRequestMetrics(moduleOwner, moduleName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch pull request metrics", moduleName)
		hasErrors = true
	}

	log.Info().Msgf("%s: uploading %d pull request metrics", moduleName, len(pullRequests))
	err = bq.Insert(config.PullRequestsTable, pullRequests)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to insert pull request metrics", moduleName)
		hasErrors = true
	}

	log.Info().Msgf("%s: fetching release metrics", moduleName)

	releases, err := metrics.GetReleaseMetrics(moduleOwner, moduleName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch release metrics", moduleName)
		hasErrors = true
	}

	log.Info().Msgf("%s: uploading %d release metrics", moduleName, len(releases))
	err = bq.Insert(config.ReleasesTable, releases)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to insert release metrics", moduleName)
		hasErrors = true
	}

	if hasErrors {
		log.Fatal().Msg("errors occurred while processing metrics you should check the previous logs for more details")
		os.Exit(1)
	} else {
		log.Info().Msg("successfully processed metrics")
	}
}
