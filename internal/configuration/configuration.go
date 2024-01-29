// Package configuration contains a number of methods that are used
// to provide configuration to the wider application. It uses viper
// to pull config from either the environment or a config file then
// unmarhsals the config into the configuration struct. The configuration struct
// is made available to the application via a package level variable
// called Config.
package configuration

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const (
	githubTokenKey          = "GITHUB_TOKEN" //nolint:gosec
	bigQueryProjectIDKey    = "BIG_QUERY_PROJECT_ID"
	bigQueryDatasetNameKey  = "BIG_QUERY_DATASET_NAME"
	issuesTableKey          = "ISSUES_TABLE"
	issuesAggTableKey       = "ISSUES_AGG_TABLE"
	pullRequestsTableKey    = "PULL_REQUESTS_TABLE"
	pullRequestsAggTableKey = "PULL_REQUESTS_AGG_TABLE"
	releasesTableKey        = "RELEASES_TABLE"
	lastRunTableKey         = "LAST_RUN_TABLE"
	repoNameKey             = "REPO_NAME"
	repoOwnerKey            = "REPO_OWNER"
)

var configMap = map[string]string{
	"GitHubToken":          "GITHUB_TOKEN",
	"BigQueryProjectID":    "BIG_QUERY_PROJECT_ID",
	"BigQueryDatasetName":  "BIG_QUERY_DATASET_NAME",
	"IssuesTable":          "ISSUES_TABLE",
	"IssuesAggTable":       "ISSUES_AGG_TABLE",
	"PullRequestsTable":    "PULL_REQUESTS_TABLE",
	"PullRequestsAggTable": "PULL_REQUESTS_AGG_TABLE",
	"ReleasesTable":        "RELEASES_TABLE",
	"LastRunTable":         "LAST_RUN_TABLE",
	"RepoName":             "REPO_NAME",
	"RepoOwner":            "REPO_OWNER",
}

var Config configuration

type configuration struct {
	GitHubToken          string `mapstructure:"github_token"`
	BigQueryProjectID    string `mapstructure:"big_query_project_id"`
	BigQueryDatasetName  string `mapstructure:"big_query_dataset_name"`
	IssuesTable          string `mapstructure:"issues_table"`
	IssuesAggTable       string `mapstructure:"issues_agg_table"`
	PullRequestsTable    string `mapstructure:"pull_requests_table"`
	PullRequestsAggTable string `mapstructure:"pull_requests_agg_table"`
	ReleasesTable        string `mapstructure:"releases_table"`
	LastRunTable         string `mapstructure:"last_run_table"`
	RepoName             string `mapstructure:"repo_name"`
	RepoOwner            string `mapstructure:"repo_owner"`
}

func InitConfig() error {
	viper.SetDefault(bigQueryDatasetNameKey, "github_metrics")
	viper.SetDefault(issuesTableKey, "issues")
	viper.SetDefault(issuesAggTableKey, "issues_agg")
	viper.SetDefault(pullRequestsTableKey, "pull_requests")
	viper.SetDefault(pullRequestsAggTableKey, "pull_requests_agg")
	viper.SetDefault(releasesTableKey, "releases")
	viper.SetDefault(lastRunTableKey, "last_run")

	_ = viper.BindEnv(githubTokenKey)
	_ = viper.BindEnv(bigQueryProjectIDKey)
	_ = viper.BindEnv(bigQueryDatasetNameKey)
	_ = viper.BindEnv(issuesTableKey)
	_ = viper.BindEnv(issuesAggTableKey)
	_ = viper.BindEnv(pullRequestsTableKey)
	_ = viper.BindEnv(pullRequestsAggTableKey)
	_ = viper.BindEnv(releasesTableKey)
	_ = viper.BindEnv(lastRunTableKey)
	_ = viper.BindEnv(repoNameKey)
	_ = viper.BindEnv(repoOwnerKey)

	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}

	return validate(Config)
}

// Needs to be way better than this..
func validate(config configuration) error {
	var missingConfig []string

	v := reflect.ValueOf(config)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			missingConfig = append(missingConfig, configMap[v.Type().Field(i).Name])
		}
	}

	if len(missingConfig) > 0 {
		return fmt.Errorf("required environment variables are missing: %s", strings.Join(missingConfig, ", "))
	}

	return nil
}
