// Package export contains the logic for exporting GitHub metrics to BigQuery.
// This is a port from the previous main.go and should be rationalized in the future.
package export

import (
	"context"
	"errors"
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

	repoName := config.RepoName
	repoOwner := config.RepoOwner

	log.Info().Msg("Initializing BigQuery client")
	ctx := context.Background()
	bq, err = bigqueryclient.NewBigQueryClient(ctx, config.BigQueryProjectID, config.BigQueryDatasetName)
	if err != nil {
		return fmt.Errorf("failed to create BigQuery client: %w", err)
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
		return errors.New("errors occurred while processing metrics you should check the previous logs for more details")
	} else {
		log.Info().Msg("successfully processed metrics")
	}

	return nil
}
