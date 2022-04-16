package main

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

type Assets []string

type Tag struct {
	Type            string
	URL             string
	Repository      string
	User            string
	Name            string
	TagName         *string
	PendingTag      *string
	Body            *string
	State           string
	Prerelease      bool
	TargetCommitish string
	ReleaseAssets   Assets
	PublishedAt     *time.Time
	CreatedAt       *time.Time
}

func getTags(project *gitlab.Project) ([]Tag, error) {
	tags, _, err := getClient().Tags.ListTags(project.ID, &gitlab.ListTagsOptions{})

	t := make([]Tag, len(tags))
	for i, tag := range tags {
		var tagName *string
		var body *string
		if tag.Release != nil {
			tagName = &tag.Release.TagName
			body = &tag.Release.Description
		}

		var authoredDate *time.Time
		if tag.Commit != nil {
			authoredDate = tag.Commit.AuthoredDate
		}

		t[i] = Tag{
			Type:            "tag",
			URL:             tag.Name,
			Repository:      project.WebURL,
			User:            "",
			Name:            project.Name,
			TagName:         tagName,
			PendingTag:      tagName,
			Body:            body,
			State:           "published",
			Prerelease:      false,
			TargetCommitish: "master",
			ReleaseAssets:   Assets{},
			PublishedAt:     authoredDate,
			CreatedAt:       authoredDate,
		}
	}

	return t, err
}
