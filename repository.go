package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Repository struct {
	Type          string     `json:"type"`
	URL           string     `json:"url"`
	Owner         string     `json:"owner"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Website       string     `json:"website"`
	Private       string     `json:"private"`
	HasIssues     bool       `json:"has_issues"`
	HasWiki       bool       `json:"has_wiki"`
	HasDownloads  bool       `json:"has_downloads"`
	Labels        []string   `json:"labels"`
	Webhooks      []string   `json:"webhooks"`
	Collaborators []string   `json:"collaborators"`
	GitUrl        string     `json:"git_url"`
	WikiUrl       string     `json:"wiki_url"`
	DefaultBranch string     `json:"default_branch"`
	CreatedAt     *time.Time `json:"created_at"`
}

func getRepository(project *gitlab.Project) {
	r, _, err := getClient().Repositories.ListTree(project.ID, &gitlab.ListTreeOptions{})
	if err != nil {
		log.Fatalf("Failed to get repository information for projectID %d: %v", project.ID, err)
	}
	fmt.Println("r", r)
}
