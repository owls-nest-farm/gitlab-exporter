package main

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

type Label struct {
	Type       string     `json:"type"`
	URL        string     `json:"url"`
	Repository string     `json:"repository"`
	User       string     `json:"user"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	Assignee   string     `json:"assignee"`
	Milestone  string     `json:"milestone"`
	Labels     []string   `json:"labels"`
	ClosedAt   *time.Time `json:"closed_at"`
	CreatedAt  *time.Time `json:"created_at"`
}

func getLabels(project *gitlab.Project) ([]*gitlab.Label, *gitlab.Response, error) {
	return getClient().Labels.ListLabels(project.ID, &gitlab.ListLabelsOptions{})
}
