package githubclient

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type Label string

type Issue struct {
	Repository string    `json:"repository"`
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Labels     []string  `json:"labels"`
	Author     string    `json:"author"`
	State      string    `json:"state"`
}

var IssueQuery struct {
	Repository struct {
		Issues struct {
			Nodes []struct {
				Number    int
				Title     string
				State     string
				CreatedAt time.Time
				UpdatedAt time.Time
				Labels    struct {
					Nodes []struct {
						Name string
					}
				} `graphql:"labels(first: 100)"`
				Author struct {
					Login string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
		} `graphql:"issues(first: 100, after: $cursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (client *githubClient) GetIssues(ctx context.Context, owner string, name string) ([]Issue, error) {
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
		"cursor": (*githubv4.String)(nil),
	}

	var metrics []Issue
	for {
		err := client.v4.Query(ctx, &IssueQuery, variables)
		if err != nil {
			return nil, err
		}

		for _, node := range IssueQuery.Repository.Issues.Nodes {
			var labels []string
			for _, label := range node.Labels.Nodes {
				labels = append(labels, label.Name)
			}

			metrics = append(metrics, Issue{
				Repository: name,
				Number:     node.Number,
				Title:      node.Title,
				CreatedAt:  node.CreatedAt,
				UpdatedAt:  node.UpdatedAt,
				Labels:     labels,
				Author:     node.Author.Login,
				State:      node.State,
			})
		}

		if !IssueQuery.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = IssueQuery.Repository.Issues.PageInfo.EndCursor
	}

	return metrics, nil
}
