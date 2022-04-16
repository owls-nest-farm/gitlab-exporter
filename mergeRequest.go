package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

type PullRequest struct {
	Type       string     `json:"type"`
	URL        string     `json:"url"`
	User       string     `json:"user"`
	Repository string     `json:"repository"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	Base       Base       `json:"base"`
	Head       Base       `json:"head"`
	Assignee   *string    `json:"assignee"`
	Milestone  *string    `json:"milestone"`
	Labels     []string   `json:"labels"`
	MergedAt   *time.Time `json:"merged_at"`
	ClosedAt   *time.Time `json:"closed_at"`
	CreatedAt  *time.Time `json:"created_at"`
}

type Base struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
	User string `json:"user"`
	Repo string `json:"repo"`
}

func getPullRequests(project *gitlab.Project) []PullRequest {
	mergeRequests, _, err := getClient().MergeRequests.ListProjectMergeRequests(project.ID, &gitlab.ListProjectMergeRequestsOptions{})
	if err != nil {
		log.Fatalf("Failed to get merge requests for projectID %d: %v", project.ID, err)
	}

	prs := make([]PullRequest, len(mergeRequests))
	for i, mergeRequest := range mergeRequests {
		var assignee *string
		if mergeRequest.Assignee != nil {
			assignee = &mergeRequest.Assignee.WebURL
		}

		var milestone *string
		if mergeRequest.Milestone != nil {
			milestone = &mergeRequest.Milestone.WebURL
		}

		var mergedAt *time.Time
		if mergeRequest.State == "merged" {
			mergedAt = mergeRequest.UpdatedAt
		}

		var closedAt *time.Time
		if mergeRequest.State == "closed" || mergeRequest.State == "merged" {
			closedAt = mergeRequest.UpdatedAt
		}

		var labels []string
		if len(mergeRequest.Labels) > 0 {
			labels = make([]string, len(mergeRequest.Labels))
			for i, label := range mergeRequest.Labels {
				labels[i] = fmt.Sprintf("%s/labels#/%s", project.WebURL, label)
			}
		}

		prs[i] = PullRequest{
			Type:       "pull_request",
			URL:        mergeRequest.WebURL,
			User:       mergeRequest.Author.WebURL,
			Repository: project.WebURL,
			Title:      mergeRequest.Title,
			Body:       mergeRequest.Description,
			Base: Base{
				Ref: mergeRequest.TargetBranch,
				//				SHA:  mergeRequest.SHA,
				//				User: project.WebURL,
				Repo: project.WebURL,
			},
			Head: Base{
				Ref: mergeRequest.SourceBranch,
				//				SHA:  mergeRequest.SHA,
				//				User: mergeRequest.Author.WebURL,
				Repo: project.WebURL,
			},
			Assignee:  assignee,
			Milestone: milestone,
			Labels:    labels,
			MergedAt:  mergedAt,
			ClosedAt:  closedAt,
			CreatedAt: mergeRequest.CreatedAt,
		}
	}

	return prs
}
