package main

import (
	"log"
	"sort"
	"time"

	"github.com/xanzy/go-gitlab"
)

type MilestoneService struct {
	*BaseService
}

type Milestone struct {
	Type        string          `json:"type"`
	URL         string          `json:"url"`
	Repository  string          `json:"repository"`
	User        string          `json:"user"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	State       string          `json:"state"`
	DueOn       *gitlab.ISOTime `json:"due_on"`
	CreatedAt   *time.Time      `json:"created_at"`
}

func NewMilestoneService(e *Exporter) *MilestoneService {
	return &MilestoneService{
		BaseService: &BaseService{
			exporter: e,
			filename: "milestones.json",
		},
	}
}

func (m *MilestoneService) GetAll() ([]Milestone, error) {
	project := m.exporter.CurrentProject

	milestones, _, err := m.exporter.Client.Milestones.ListMilestones(project.ID, &gitlab.ListMilestonesOptions{})

	ms := make([]Milestone, len(milestones))
	for i, milestone := range milestones {
		state := "open"
		if milestone.State != "active" {
			state = "closed"
		}

		ms[i] = Milestone{
			Type:        "milestone",
			URL:         milestone.WebURL,
			Repository:  project.WebURL,
			User:        project.WebURL,
			Title:       milestone.Title,
			Description: milestone.Description,
			State:       state,
			DueOn:       milestone.DueDate,
			CreatedAt:   milestone.CreatedAt,
		}
	}

	return ms, err
}

func (m *MilestoneService) Export() {
	milestones, err := m.GetAll()
	if err != nil {
		log.Fatalf("Failed to get milestones for projectID %d: %v", m.exporter.CurrentProject.ID, err)
	}
	if len(milestones) > 0 {
		m.exporter.State.Milestones = append(m.exporter.State.Milestones, milestones...)
	}
}

func (m *MilestoneService) WriteFile() error {
	milestones := m.exporter.State.MergeRequests
	sort.Slice(milestones, func(i, j int) bool {
		return milestones[i].Title > milestones[j].Title
	})
	return m.exporter.WriteJsonFile(m.filename, milestones)
}
