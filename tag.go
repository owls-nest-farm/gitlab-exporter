package main

import (
	"sort"
	"time"

	"github.com/xanzy/go-gitlab"
)

type TagService struct {
	*BaseService
}

type Assets []string

type Tag struct {
	Type            string     `json:"type"`
	URL             string     `json:"url"`
	Repository      string     `json:"repository"`
	User            string     `json:"user"`
	Name            string     `json:"name"`
	TagName         *string    `json:"tag_name"`
	PendingTag      *string    `json:"pending_tag"`
	Body            *string    `json:"body"`
	State           string     `json:"state"`
	Prerelease      bool       `json:"prerelease"`
	TargetCommitish string     `json:"target_commitish"`
	ReleaseAssets   Assets     `json:"release_assets"`
	PublishedAt     *time.Time `json:"published_at"`
	CreatedAt       *time.Time `json:"created_at"`
}

func NewTagService(e *Exporter) *TagService {
	return &TagService{
		BaseService: &BaseService{
			exporter: e,
			filename: "releases.json",
		},
	}
}

func (t *TagService) GetAll() ([]Tag, error) {
	project := t.exporter.CurrentProject

	tags, _, err := t.exporter.Client.Tags.ListTags(project.ID, &gitlab.ListTagsOptions{})

	ts := make([]Tag, len(tags))
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

		ts[i] = Tag{
			Type:            "release",
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

	return ts, err
}

func (t *TagService) Export() {
	tags, err := t.GetAll()
	if err != nil {
		panic(err)
	}
	if len(tags) > 0 {
		t.exporter.State.Tags = append(t.exporter.State.Tags, tags...)
	}
}

func (t *TagService) WriteFile() error {
	tags := t.exporter.State.Tags
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name > tags[j].Name
	})
	return t.exporter.WriteJsonFile(t.filename, tags)
}
