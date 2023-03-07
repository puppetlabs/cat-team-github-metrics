package configuration

import (
	"testing"
)

func ExpectEqual(t *testing.T, check, expected interface{}) bool {
	if check != expected {
		t.Helper()
		t.Errorf("Expected %s, actual %s", expected, check)
		return false
	}
	return true
}

func Test_Validate_Function(t *testing.T) {
	t.Log("When all the attributes for Config struct is passed to validate method")
	{
		t.Run("it should append no element in the array", func(t *testing.T) {
			var config = configuration{
				GitHubToken:          "test-githubTokenKey",
				BigQueryProjectID:    "test-bigQueryProjectIDKey",
				BigQueryDatasetName:  "test-bigQueryDatasetNameKey",
				IssuesTable:          "test-issuesTableKey",
				IssuesAggTable:       "test-issuesAggTableKey",
				PullRequestsTable:    "test-pullRequestsTableKey",
				PullRequestsAggTable: "test-pullRequestsAggTableKey",
				ReleasesTable:        "test-releasesTableKey",
				LastRunTable:         "test-lastRunTableKey",
				RepoName:             "test-repoNameKey",
				RepoOwner:            "test-repoOwner",
			}
			testValidate := validate(config)
			ExpectEqual(t, testValidate, nil)
		})
	}
	t.Log("When RepoOwner for Config struct is passed to validate method")
	{
		t.Run("it should return error", func(t *testing.T) {
			var config = configuration{
				GitHubToken:          "test-githubTokenKey",
				BigQueryProjectID:    "test-bigQueryProjectIDKey",
				BigQueryDatasetName:  "test-bigQueryDatasetNameKey",
				IssuesTable:          "test-issuesTableKey",
				IssuesAggTable:       "test-issuesAggTableKey",
				PullRequestsTable:    "test-pullRequestsTableKey",
				PullRequestsAggTable: "test-pullRequestsAggTableKey",
				ReleasesTable:        "test-releasesTableKey",
				LastRunTable:         "test-lastRunTableKey",
				RepoName:             "",
			}
			testValidate := validate(config)
			ExpectEqual(t, testValidate.Error(), "required environment variables are missing: REPO_NAME, REPO_OWNER")
		})
	}
}
