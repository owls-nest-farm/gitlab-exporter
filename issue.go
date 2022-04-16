package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Issue struct {
	Type       string     `json:"type"`
	URL        string     `json:"url"`
	Repository string     `json:"repository"`
	User       string     `json:"user"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	Assignee   *string    `json:"assignee"`
	Milestone  *string    `json:"milestone"`
	Labels     []string   `json:"labels"`
	ClosedAt   *time.Time `json:"closed_at"`
	CreatedAt  *time.Time `json:"created_at"`
}

func getIssues(project *gitlab.Project) []Issue {
	issues, _, err := getClient().Issues.ListProjectIssues(project.ID, &gitlab.ListProjectIssuesOptions{})
	if err != nil {
		log.Fatalf("Failed to get issues for projectID %d: %v", project.ID, err)
	}

	i := make([]Issue, len(issues))
	for j, issue := range issues {
		var assignee *string
		if issue.Assignee != nil {
			assignee = &issue.Assignee.WebURL
		}

		var milestone *string
		if issue.Milestone != nil {
			milestone = &issue.Milestone.WebURL
		}

		var closedAt *time.Time
		if issue.State == "closed" {
			closedAt = issue.UpdatedAt
		}

		var labels []string
		if len(issue.Labels) > 0 {
			labels = make([]string, len(issue.Labels))
			for i, label := range issue.Labels {
				labels[i] = fmt.Sprintf("%s/labels#/%s", project.WebURL, label)
			}
		}

		i[j] = Issue{
			Type:       "issue",
			URL:        issue.WebURL,
			Repository: project.WebURL,
			User:       issue.Author.WebURL,
			Title:      issue.Title,
			Body:       issue.Description,
			Assignee:   assignee,
			Milestone:  milestone,
			Labels:     labels,
			ClosedAt:   closedAt,
			CreatedAt:  issue.CreatedAt,
		}
	}

	return i
}
