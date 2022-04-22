package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/xanzy/go-gitlab"
)

type RepositoryService struct {
	*BaseService
}

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

func NewRepositoryService(e *Exporter) *RepositoryService {
	return &RepositoryService{
		BaseService: &BaseService{
			exporter: e,
			filename: "repositories.json",
		},
	}
}

func (r *RepositoryService) Clone(project *gitlab.Project) (*git.Repository, error) {
	// Will be something like `repositories/gl-group1/gl-subgroup1/quux.git`.
	dir := fmt.Sprintf("%s/%s.git", r.exporter.TmpRepositoryDir, project.PathWithNamespace)
	return git.PlainClone(dir, true, &git.CloneOptions{
		URL: project.HTTPURLToRepo,
	})
}

func (r *RepositoryService) Get() {
	project := r.exporter.CurrentProject

	repos, _, err := r.exporter.Client.Repositories.ListTree(project.ID, &gitlab.ListTreeOptions{})
	if err != nil {
		log.Fatalf("Failed to get repository information for projectID %d: %v", project.ID, err)
	}
	fmt.Println("repos", repos)
}
