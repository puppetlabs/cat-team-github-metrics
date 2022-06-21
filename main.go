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

	log.Info().Msg("Initializing BigQuery client")
	ctx := context.Background()
	bq, err = bigqueryclient.NewBigQueryClient(ctx, config.BigQueryProjectID, config.BigQueryDatasetName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create BigQuery client")
		os.Exit(1)
	}

	log.Info().Msg("fetching supported modules")
	modulesClient := modules.NewModuleClient(nil, "")
	suppportedModules, err := modulesClient.GetSupportedModules(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get supported modules")
		os.Exit(1)
	}

	log.Info().Msgf("there are %d modules to process", len(*suppportedModules))

	var hasErrors bool
	log.Info().Msg("starting to collect metrics")
	for _, module := range *suppportedModules {
		log.Info().Msgf("%s: fetching issue metrics", module.Name)
		issues, err := metrics.GetIssueMetrics(module.Owner, module.Name)
		if err != nil {
			log.Error().Err(err).Msgf("%s: failed to fetch issue metrics", module.Name)
			hasErrors = true
			continue
		}

		log.Info().Msgf("%s: uploading %d issue(s) metrics", module.Name, len(issues))
		err = bq.Insert(config.IssuesTable, issues)
		if err != nil {
			log.Error().Err(err).Msgf("%s: failed to insert issue metrics", module.Name)
			hasErrors = true
			continue
		}

		log.Info().Msgf("%s: fetching pull request metrics", module.Name)
		pullRequests, err := metrics.GetPullRequestMetrics(module.Owner, module.Name)
		if err != nil {
			log.Error().Err(err).Msgf("%s: failed to fetch pull request metrics", module.Name)
			hasErrors = true
			continue
		}

		log.Info().Msgf("%s: uploading %d pull request(s) metrics", module.Name, len(pullRequests))
		err = bq.Insert(config.PullRequestsTable, pullRequests)
		if err != nil {
			log.Error().Err(err).Msgf("%s: failed to insert pull request metrics", module.Name)
			hasErrors = true
			continue
		}
	}

	if hasErrors {
		log.Fatal().Msg("errors occurred while processing metrics you should check the previous logs for more details")
		os.Exit(1)
	} else {
		log.Info().Msg("successfully processed metrics")
	}
}

// func setup() {
// 	var err error
//
// 	fmt.Println("Creating issues table")
// 	err = bq.CreateTable(issuesTable, metrics.IssueMetric{})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	fmt.Println("Creating pull requests table")
// 	err = bq.CreateTable(pullRequestsTable, metrics.PullRequestMetric{})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
