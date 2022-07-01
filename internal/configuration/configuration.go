//Package configuration contains a number of methods that are used
//to provide configuration to the wider application. It uses viper
//to pull config from either the environment or a config file then
//unmarhsals the config into the configuration struct. The configuration struct
//is made available to the application via a package level variable
//called Config.
package configuration

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	githubTokenKey         = "GITHUB_TOKEN"
	bigQueryProjectIDKey   = "BIG_QUERY_PROJECT_ID"
	bigQueryDatasetNameKey = "BIG_QUERY_DATASET_NAME"
	issuesTableKey         = "ISSUES_TABLE"
	pullRequestsTableKey   = "PULL_REQUESTS_TABLE"
	releasesTableKey       = "RELEASES_TABLE"
	moduleNameKey          = "MODULE_NAME"
	moduleOwnerKey         = "MODULE_OWNER"
)

var Config configuration

type configuration struct {
	GitHubToken         string `mapstructure:"github_token"`
	BigQueryProjectID   string `mapstructure:"big_query_project_id"`
	BigQueryDatasetName string `mapstructure:"big_query_dataset_name"`
	IssuesTable         string `mapstructure:"issues_table"`
	PullRequestsTable   string `mapstructure:"pull_requests_table"`
	ReleasesTable       string `mapstructure:"releases_table"`
	ModuleName          string `mapstructure:"module_name"`
	ModuleOwner         string `mapstructure:"module_owner"`
}

func InitConfig() error {
	viper.SetDefault(bigQueryDatasetNameKey, "github_metrics")
	viper.SetDefault(issuesTableKey, "issues")
	viper.SetDefault(pullRequestsTableKey, "pull_requests")
	viper.SetDefault(releasesTableKey, "releases")

	_ = viper.BindEnv(githubTokenKey)
	_ = viper.BindEnv(bigQueryProjectIDKey)
	_ = viper.BindEnv(bigQueryDatasetNameKey)
	_ = viper.BindEnv(issuesTableKey)
	_ = viper.BindEnv(pullRequestsTableKey)
	_ = viper.BindEnv(releasesTableKey)
	_ = viper.BindEnv(moduleNameKey)
	_ = viper.BindEnv(moduleOwnerKey)

	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}

	return validate(Config)
}

func validate(config configuration) error {
	var missingConfig []string

	if config.GitHubToken == "" {
		missingConfig = append(missingConfig, githubTokenKey)
	}

	if config.BigQueryProjectID == "" {
		missingConfig = append(missingConfig, bigQueryProjectIDKey)
	}

	if config.BigQueryDatasetName == "" {
		missingConfig = append(missingConfig, bigQueryDatasetNameKey)
	}

	if config.IssuesTable == "" {
		missingConfig = append(missingConfig, issuesTableKey)
	}

	if config.PullRequestsTable == "" {
		missingConfig = append(missingConfig, pullRequestsTableKey)
	}

	if config.ModuleName == "" {
		missingConfig = append(missingConfig, moduleNameKey)
	}

	if config.ModuleOwner == "" {
		missingConfig = append(missingConfig, moduleOwnerKey)
	}

	if len(missingConfig) > 0 {
		return fmt.Errorf("required environment variables are missing: %s", strings.Join(missingConfig, ", "))
	}

	return nil
}
