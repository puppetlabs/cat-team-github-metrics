package githubclient

import (
	"context"

	"github.com/shurcooL/githubv4"
)

type IssueCount struct {
	Repository string `json:"repository"`
	Count      int    `json:"count"`
}

var OpenIssueCountQuery struct {
	Repository struct {
		Issues struct {
			TotalCount int
		} `graphql:"issues(states: OPEN)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (client *githubClient) GetOpenIssueCount(ctx context.Context, owner string, name string) (IssueCount, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := client.v4.Query(ctx, &OpenIssueCountQuery, variables)
	if err != nil {
		return IssueCount{}, err
	}

	openIssues := IssueCount{
		Repository: name,
		Count:      OpenIssueCountQuery.Repository.Issues.TotalCount,
	}

	return openIssues, nil
}
