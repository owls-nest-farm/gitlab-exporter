package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

type MergeRequestService struct {
	*BaseService
}

type MergeRequest struct {
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

func NewMergeRequestService(e *Exporter) *MergeRequestService {
	field := "Description"
	return &MergeRequestService{
		BaseService: &BaseService{
			exporter:        e,
			filename:        "merge_requests.json",
			attachmentField: &field,
		},
	}
}

func (m *MergeRequestService) GetAll() ([]MergeRequest, error) {
	project := m.exporter.CurrentProject

	mergeRequests, _, err := m.exporter.Client.MergeRequests.ListProjectMergeRequests(project.ID, &gitlab.ListProjectMergeRequestsOptions{})
	if err != nil {
		return nil, err
	}

	mrs := make([]MergeRequest, len(mergeRequests))
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

		if m.attachmentField != nil {
			m.exporter.Attachments.Export(*m.attachmentField, mergeRequest)
		}

		mrs[i] = MergeRequest{
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

	return mrs, nil
}

func (m *MergeRequestService) Export() {
	mergeRequests, err := m.GetAll()
	if err != nil {
		log.Fatalf("Failed to get merge requests for projectID %d: %v", m.exporter.CurrentProject.ID, err)
	}
	if len(mergeRequests) > 0 {
		m.exporter.State.MergeRequests = append(m.exporter.State.MergeRequests, mergeRequests...)
	}
}

func (m *MergeRequestService) WriteFile() error {
	return m.exporter.WriteJsonFile(m.filename, m.exporter.State.MergeRequests)
}
