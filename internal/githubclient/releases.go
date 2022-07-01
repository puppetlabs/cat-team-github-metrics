package githubclient

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type Release struct {
	Repository  string
	Name        string
	PublishedAt time.Time
}

var ReleaseQuery struct {
	Repository struct {
		LatestRelease struct {
			TagName     string
			PublishedAt time.Time
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (client *githubClient) GetLatestRelease(ctx context.Context, owner string, name string) ([]Release, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := client.v4.Query(ctx, &ReleaseQuery, variables)
	if err != nil {
		return nil, err
	}

	var metric []Release

	metric = append(metric, Release{
		Repository:  name,
		Name:        ReleaseQuery.Repository.LatestRelease.TagName,
		PublishedAt: ReleaseQuery.Repository.LatestRelease.PublishedAt,
	})

	return metric, nil
}
