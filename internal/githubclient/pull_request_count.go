package githubclient

import (
	"context"

	"github.com/shurcooL/githubv4"
)

type PullRequestCount struct {
	Repository string `json:"repository"`
	Count      int    `json:"count"`
}

var OpenPullRequestCountQuery struct {
	Repository struct {
		PullRequests struct {
			TotalCount int
		} `graphql:"pullRequests(states: OPEN)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (client *githubClient) GetOpenPullRequestCount(ctx context.Context, owner string, name string) (PullRequestCount, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := client.v4.Query(ctx, &OpenPullRequestCountQuery, variables)
	if err != nil {
		return PullRequestCount{}, err
	}

	openPullRequests := PullRequestCount{
		Repository: name,
		Count:      OpenPullRequestCountQuery.Repository.PullRequests.TotalCount,
	}

	return openPullRequests, nil
}
