package metrics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/puppetlabs/cat-team-github-metrics/internal/githubclient"
)

// mockGitHubClient implements githubclient.GitHubClient for testing.
type mockGitHubClient struct {
	issues           []githubclient.Issue
	issueCount       githubclient.IssueCount
	pullRequests     []githubclient.PullRequest
	pullRequestCount githubclient.PullRequestCount
	releases         []githubclient.Release
	err              error
}

func (m *mockGitHubClient) GetIssues(_ context.Context, _, _ string) ([]githubclient.Issue, error) {
	return m.issues, m.err
}

func (m *mockGitHubClient) GetOpenIssueCount(_ context.Context, _, _ string) (githubclient.IssueCount, error) {
	return m.issueCount, m.err
}

func (m *mockGitHubClient) GetPullRequests(_ context.Context, _, _ string) ([]githubclient.PullRequest, error) {
	return m.pullRequests, m.err
}

func (m *mockGitHubClient) GetOpenPullRequestCount(_ context.Context, _, _ string) (githubclient.PullRequestCount, error) {
	return m.pullRequestCount, m.err
}

func (m *mockGitHubClient) GetLatestRelease(_ context.Context, _, _ string) ([]githubclient.Release, error) {
	return m.releases, m.err
}

// helpers

var (
	testTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	errAPI   = errors.New("api error")
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- IssueMetric.Save ---

func TestIssueMetricSave(t *testing.T) {
	m := &IssueMetric{
		Repository:     "my-repo",
		Number:         7,
		Title:          "a bug",
		CreatedAt:      testTime,
		UpdatedAt:      testTime,
		Labels:         []string{"bug"},
		Author:         "alice",
		State:          "OPEN",
		CollectionTime: testTime,
	}

	row, insertID, err := m.Save()
	assertNoError(t, err)
	if insertID != "" {
		t.Errorf("expected empty insertID, got %q", insertID)
	}
	if row["Repository"] != "my-repo" {
		t.Errorf("Repository: got %v", row["Repository"])
	}
	if row["Number"] != 7 {
		t.Errorf("Number: got %v", row["Number"])
	}
	if row["Title"] != "a bug" {
		t.Errorf("Title: got %v", row["Title"])
	}
	if row["CreatedAt"] != testTime.Unix() {
		t.Errorf("CreatedAt: got %v", row["CreatedAt"])
	}
	if row["UpdatedAt"] != testTime.Unix() {
		t.Errorf("UpdatedAt: got %v", row["UpdatedAt"])
	}
	if row["Author"] != "alice" {
		t.Errorf("Author: got %v", row["Author"])
	}
	if row["State"] != "OPEN" {
		t.Errorf("State: got %v", row["State"])
	}
	if row["CollectionTime"] != testTime.Unix() {
		t.Errorf("CollectionTime: got %v", row["CollectionTime"])
	}
}

// --- IssueAggregatedMetric.Save ---

func TestIssueAggregatedMetricSave(t *testing.T) {
	m := &IssueAggregatedMetric{
		Repository:     "my-repo",
		Count:          42,
		CollectionTime: testTime,
	}

	row, insertID, err := m.Save()
	assertNoError(t, err)
	if insertID != "" {
		t.Errorf("expected empty insertID, got %q", insertID)
	}
	if row["Repository"] != "my-repo" {
		t.Errorf("Repository: got %v", row["Repository"])
	}
	if row["Count"] != 42 {
		t.Errorf("Count: got %v", row["Count"])
	}
	if row["CollectionTime"] != testTime.Unix() {
		t.Errorf("CollectionTime: got %v", row["CollectionTime"])
	}
}

// --- PullRequestMetric.Save ---

func TestPullRequestMetricSave(t *testing.T) {
	m := &PullRequestMetric{
		Repository:     "my-repo",
		Number:         3,
		Title:          "a PR",
		CreatedAt:      testTime,
		UpdatedAt:      testTime,
		Labels:         []string{"enhancement"},
		Author:         "bob",
		State:          "MERGED",
		Merged:         true,
		CollectionTime: testTime,
	}

	row, insertID, err := m.Save()
	assertNoError(t, err)
	if insertID != "" {
		t.Errorf("expected empty insertID, got %q", insertID)
	}
	if row["Repository"] != "my-repo" {
		t.Errorf("Repository: got %v", row["Repository"])
	}
	if row["Number"] != 3 {
		t.Errorf("Number: got %v", row["Number"])
	}
	if row["Merged"] != true {
		t.Errorf("Merged: got %v", row["Merged"])
	}
	if row["CollectionTime"] != testTime.Unix() {
		t.Errorf("CollectionTime: got %v", row["CollectionTime"])
	}
}

// --- PullRequestAggregatedMetric.Save ---

func TestPullRequestAggregatedMetricSave(t *testing.T) {
	m := &PullRequestAggregatedMetric{
		Repository:     "my-repo",
		Count:          5,
		CollectionTime: testTime,
	}

	row, insertID, err := m.Save()
	assertNoError(t, err)
	if insertID != "" {
		t.Errorf("expected empty insertID, got %q", insertID)
	}
	if row["Repository"] != "my-repo" {
		t.Errorf("Repository: got %v", row["Repository"])
	}
	if row["Count"] != 5 {
		t.Errorf("Count: got %v", row["Count"])
	}
	if row["CollectionTime"] != testTime.Unix() {
		t.Errorf("CollectionTime: got %v", row["CollectionTime"])
	}
}

// --- ReleaseMetric.Save ---

func TestReleaseMetricSave(t *testing.T) {
	m := &ReleaseMetric{
		Repository:     "my-repo",
		Name:           "v1.2.3",
		PublishedAt:    testTime,
		CollectionTime: testTime,
	}

	row, insertID, err := m.Save()
	assertNoError(t, err)
	if insertID != "" {
		t.Errorf("expected empty insertID, got %q", insertID)
	}
	if row["Repository"] != "my-repo" {
		t.Errorf("Repository: got %v", row["Repository"])
	}
	if row["Name"] != "v1.2.3" {
		t.Errorf("Name: got %v", row["Name"])
	}
	if row["PublishedAt"] != testTime.Unix() {
		t.Errorf("PublishedAt: got %v", row["PublishedAt"])
	}
	if row["CollectionTime"] != testTime.Unix() {
		t.Errorf("CollectionTime: got %v", row["CollectionTime"])
	}
}

// --- LastRunMetric.Save ---

func TestLastRunMetricSave(t *testing.T) {
	m := &LastRunMetric{
		LastRunTime:    testTime,
		CollectionTime: testTime,
	}

	row, insertID, err := m.Save()
	assertNoError(t, err)
	if insertID != "" {
		t.Errorf("expected empty insertID, got %q", insertID)
	}
	if row["LastRunTime"] != testTime.Unix() {
		t.Errorf("LastRunTime: got %v", row["LastRunTime"])
	}
	if row["CollectionTime"] != testTime.Unix() {
		t.Errorf("CollectionTime: got %v", row["CollectionTime"])
	}
}

// --- mapIssueMetrics ---

func TestMapIssueMetrics(t *testing.T) {
	input := []githubclient.Issue{
		{
			Repository: "repo-a",
			Number:     1,
			Title:      "issue one",
			CreatedAt:  testTime,
			UpdatedAt:  testTime,
			Labels:     []string{"bug"},
			Author:     "alice",
			State:      "OPEN",
		},
	}

	result := mapIssueMetrics(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Repository != "repo-a" {
		t.Errorf("Repository: got %v", result[0].Repository)
	}
	if result[0].Number != 1 {
		t.Errorf("Number: got %v", result[0].Number)
	}
	if result[0].Title != "issue one" {
		t.Errorf("Title: got %v", result[0].Title)
	}
	if result[0].Author != "alice" {
		t.Errorf("Author: got %v", result[0].Author)
	}
	if result[0].State != "OPEN" {
		t.Errorf("State: got %v", result[0].State)
	}
}

func TestMapIssueMetricsEmpty(t *testing.T) {
	result := mapIssueMetrics([]githubclient.Issue{})
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

// --- mapPullRequestMetrics ---

func TestMapPullRequestMetrics(t *testing.T) {
	input := []githubclient.PullRequest{
		{
			Repository: "repo-a",
			Number:     2,
			Title:      "pr one",
			CreatedAt:  testTime,
			UpdatedAt:  testTime,
			Labels:     []string{"enhancement"},
			Author:     "bob",
			State:      "MERGED",
			Merged:     true,
		},
	}

	result := mapPullRequestMetrics(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Repository != "repo-a" {
		t.Errorf("Repository: got %v", result[0].Repository)
	}
	if result[0].Merged != true {
		t.Errorf("Merged: got %v", result[0].Merged)
	}
	if result[0].State != "MERGED" {
		t.Errorf("State: got %v", result[0].State)
	}
}

func TestMapPullRequestMetricsEmpty(t *testing.T) {
	result := mapPullRequestMetrics([]githubclient.PullRequest{})
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

// --- mapReleaseMetrics ---

func TestMapReleaseMetrics(t *testing.T) {
	input := []githubclient.Release{
		{
			Repository:  "repo-a",
			Name:        "v2.0.0",
			PublishedAt: testTime,
		},
	}

	result := mapReleaseMetrics(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Repository != "repo-a" {
		t.Errorf("Repository: got %v", result[0].Repository)
	}
	if result[0].Name != "v2.0.0" {
		t.Errorf("Name: got %v", result[0].Name)
	}
	if result[0].PublishedAt != testTime {
		t.Errorf("PublishedAt: got %v", result[0].PublishedAt)
	}
}

func TestMapReleaseMetricsEmpty(t *testing.T) {
	result := mapReleaseMetrics([]githubclient.Release{})
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

// --- exported Get* functions: no-token error path ---

func TestGetIssueMetrics_NoToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	_, err := GetIssueMetrics("org", "repo")
	assertError(t, err)
}

func TestGetIssueAggregatedMetrics_NoToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	_, err := GetIssueAggregatedMetrics("org", "repo")
	assertError(t, err)
}

func TestGetPullRequestMetrics_NoToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	_, err := GetPullRequestMetrics("org", "repo")
	assertError(t, err)
}

func TestGetPullRequestAggregatedMetrics_NoToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	_, err := GetPullRequestAggregatedMetrics("org", "repo")
	assertError(t, err)
}

func TestGetReleaseMetrics_NoToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	_, err := GetReleaseMetrics("org", "repo")
	assertError(t, err)
}

// --- helper functions via mock: success paths ---

func TestGetIssueMetrics_Success(t *testing.T) {
	client := &mockGitHubClient{
		issues: []githubclient.Issue{
			{Repository: "repo", Number: 1, Title: "t", Author: "a", State: "OPEN"},
		},
	}

	result, err := getIssueMetrics(client, "org", "repo")
	assertNoError(t, err)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Number != 1 {
		t.Errorf("Number: got %v", result[0].Number)
	}
}

func TestGetIssueMetrics_Error(t *testing.T) {
	client := &mockGitHubClient{err: errAPI}
	_, err := getIssueMetrics(client, "org", "repo")
	assertError(t, err)
}

func TestGetIssueAggregatedMetrics_Success(t *testing.T) {
	client := &mockGitHubClient{
		issueCount: githubclient.IssueCount{Repository: "repo", Count: 10},
	}

	result, err := getIssueAggregatedMetrics(client, "org", "repo")
	assertNoError(t, err)
	if result.Repository != "repo" {
		t.Errorf("Repository: got %v", result.Repository)
	}
	if result.Count != 10 {
		t.Errorf("Count: got %v", result.Count)
	}
}

func TestGetIssueAggregatedMetrics_Error(t *testing.T) {
	client := &mockGitHubClient{err: errAPI}
	_, err := getIssueAggregatedMetrics(client, "org", "repo")
	assertError(t, err)
}

func TestGetPullRequestMetrics_Success(t *testing.T) {
	client := &mockGitHubClient{
		pullRequests: []githubclient.PullRequest{
			{Repository: "repo", Number: 5, Merged: true},
		},
	}

	result, err := getPullRequestMetrics(client, "org", "repo")
	assertNoError(t, err)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Merged != true {
		t.Errorf("Merged: got %v", result[0].Merged)
	}
}

func TestGetPullRequestMetrics_Error(t *testing.T) {
	client := &mockGitHubClient{err: errAPI}
	_, err := getPullRequestMetrics(client, "org", "repo")
	assertError(t, err)
}

func TestGetPullRequestAggregatedMetrics_Success(t *testing.T) {
	client := &mockGitHubClient{
		pullRequestCount: githubclient.PullRequestCount{Repository: "repo", Count: 3},
	}

	result, err := getPullRequestAggregatedMetrics(client, "org", "repo")
	assertNoError(t, err)
	if result.Count != 3 {
		t.Errorf("Count: got %v", result.Count)
	}
}

func TestGetPullRequestAggregatedMetrics_Error(t *testing.T) {
	client := &mockGitHubClient{err: errAPI}
	_, err := getPullRequestAggregatedMetrics(client, "org", "repo")
	assertError(t, err)
}

func TestGetReleaseMetrics_Success(t *testing.T) {
	client := &mockGitHubClient{
		releases: []githubclient.Release{
			{Repository: "repo", Name: "v1.0.0", PublishedAt: testTime},
		},
	}

	result, err := getReleaseMetrics(client, "org", "repo")
	assertNoError(t, err)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Name != "v1.0.0" {
		t.Errorf("Name: got %v", result[0].Name)
	}
}

func TestGetReleaseMetrics_Error(t *testing.T) {
	client := &mockGitHubClient{err: errAPI}
	_, err := getReleaseMetrics(client, "org", "repo")
	assertError(t, err)
}
