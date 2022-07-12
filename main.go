package main

import (
	"context"
	"os"
	"time"

	"github.com/puppetlabs/cat-team-github-metrics/internal/bigqueryclient"
	"github.com/puppetlabs/cat-team-github-metrics/internal/configuration"
	"github.com/puppetlabs/cat-team-github-metrics/internal/metrics"
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

	repoName := config.RepoName
	repoOwner := config.RepoOwner

	log.Info().Msg("Initializing BigQuery client")
	ctx := context.Background()
	bq, err = bigqueryclient.NewBigQueryClient(ctx, config.BigQueryProjectID, config.BigQueryDatasetName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create BigQuery client")
		os.Exit(1)
	}

	var hasErrors bool

	log.Info().Msg("starting to collect metrics")
	log.Info().Msgf("%s: fetching issue metrics", repoName)

	issues, err := metrics.GetIssueMetrics(repoOwner, repoName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch issue metrics", repoName)
		hasErrors = true
	}

	log.Info().Msgf("%s: uploading %d issue metrics", repoName, len(issues))
	err = bq.Insert(config.IssuesTable, issues)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to insert issue metrics", repoName)
		hasErrors = true
	}

	log.Info().Msgf("%s: fetching pull request metrics", repoName)
	pullRequests, err := metrics.GetPullRequestMetrics(repoOwner, repoName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch pull request metrics", repoName)
		hasErrors = true
	}

	log.Info().Msgf("%s: uploading %d pull request metrics", repoName, len(pullRequests))
	err = bq.Insert(config.PullRequestsTable, pullRequests)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to insert pull request metrics", repoName)
		hasErrors = true
	}

	log.Info().Msgf("%s: fetching release metrics", repoName)

	releases, err := metrics.GetReleaseMetrics(repoOwner, repoName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch release metrics", repoName)
		hasErrors = true
	}

	log.Info().Msgf("%s: uploading %d release metrics", repoName, len(releases))
	err = bq.Insert(config.ReleasesTable, releases)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to insert release metrics", repoName)
		hasErrors = true
	}

	if hasErrors {
		log.Fatal().Msg("errors occurred while processing metrics you should check the previous logs for more details")
		os.Exit(1)
	} else {
		log.Info().Msg("successfully processed metrics")
	}
}
