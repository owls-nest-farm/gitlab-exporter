package main

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

type Milestone struct {
	Type        string
	URL         string
	Repository  string
	User        string
	Title       string
	Description string
	State       string
	DueOn       *gitlab.ISOTime
	CreatedAt   *time.Time
}

func getMilestones(project *gitlab.Project) ([]Milestone, error) {
	milestones, _, err := getClient().Milestones.ListMilestones(project.ID, &gitlab.ListMilestonesOptions{})

	m := make([]Milestone, len(milestones))
	for i, milestone := range milestones {
		state := "open"
		if milestone.State != "active" {
			state = "closed"
		}

		m[i] = Milestone{
			Type:        "milestone",
			URL:         milestone.WebURL,
			Repository:  project.WebURL,
			User:        "",
			Title:       milestone.Title,
			Description: milestone.Description,
			State:       state,
			DueOn:       milestone.DueDate,
			CreatedAt:   milestone.CreatedAt,
		}
	}

	return m, err
}
