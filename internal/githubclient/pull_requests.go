package githubclient

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/net/context"
)

type PullRequest struct {
	Repository string    `json:"repository"`
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	State      string    `json:"state"`
	Merged     bool      `json:"merged"`
	Labels     []string  `json:"labels"`
	Author     string    `json:"author"`
}

var PRQuery struct {
	Repository struct {
		PullRequests struct {
			Nodes []struct {
				Number    int
				Title     string
				State     string
				CreatedAt time.Time
				UpdatedAt time.Time
				IsDraft   bool
				Labels    struct {
					Nodes []struct {
						Name string
					}
				} `graphql:"labels(first: 100)"`
				Author struct {
					Login string
				}
				Merged bool
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
		} `graphql:"pullRequests(first: 100, after: $cursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (client *githubClient) GetPullRequests(ctx context.Context, owner string, name string) ([]PullRequest, error) {
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
		"cursor": (*githubv4.String)(nil),
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	var metrics []PullRequest
	for {
		err := client.v4.Query(ctx, &PRQuery, variables)
		if err != nil {
			return nil, err
		}

		for _, node := range PRQuery.Repository.PullRequests.Nodes {
			var labels []string
			for _, label := range node.Labels.Nodes {
				labels = append(labels, label.Name)
			}

			if node.IsDraft && node.State == "OPEN" {
				log.Info().Msgf("Skipping draft PR %s/%d", name, node.Number)
				continue
			}

			metrics = append(metrics, PullRequest{
				Repository: name,
				Number:     node.Number,
				Title:      node.Title,
				CreatedAt:  node.CreatedAt,
				UpdatedAt:  node.UpdatedAt,
				Labels:     labels,
				Author:     node.Author.Login,
				State:      node.State,
				Merged:     node.Merged,
			})
		}

		if !PRQuery.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = PRQuery.Repository.PullRequests.PageInfo.EndCursor
	}

	return metrics, nil
}
