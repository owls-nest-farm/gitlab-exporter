package main

import (
	"fmt"
	"sort"
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
	Website       *string    `json:"website"`
	Private       bool       `json:"private"`
	HasIssues     bool       `json:"has_issues"`
	HasWiki       bool       `json:"has_wiki"`
	HasDownloads  bool       `json:"has_downloads"`
	Labels        []Label    `json:"labels"`
	Webhooks      []Webhook  `json:"webhooks"`
	Collaborators []User     `json:"collaborators"`
	GitURL        string     `json:"git_url"`
	WikiURL       *string    `json:"wiki_url" omitempty:"wiki_url"`
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

//func (r *RepositoryService) GetAll(gid string) ([]*gitlab.TreeNode, *gitlab.Response, error) {
//	return r.exporter.Client.Repositories.ListTree(gid, &gitlab.ListTreeOptions{})
//}

func (r *RepositoryService) Export(project *gitlab.Project) {
	state := r.exporter.State
	var private bool
	if project.Visibility == "private" {
		private = true
	} else {
		private = false
	}
	var owner string
	if project.Namespace != nil {
		owner = project.Namespace.WebURL
	}
	repo := Repository{
		Type:          "repository",
		URL:           project.WebURL,
		Name:          project.Name,
		Description:   project.Description,
		Owner:         owner,
		DefaultBranch: project.DefaultBranch,
		Private:       private,
		HasDownloads:  project.JobsEnabled,
		HasIssues:     project.IssuesEnabled,
		HasWiki:       project.WikiEnabled,
		Website:       nil,
		CreatedAt:     project.CreatedAt,
		//		Collaborators: state.Users,
		Collaborators: []User{},
		Labels:        state.Labels,
		Webhooks:      state.Webhooks,
		GitURL:        fmt.Sprintf("tarball://root/repositories/%s.git", project.PathWithNamespace),
	}
	var wikiURL string
	if project.WikiEnabled {
		wikiURL = fmt.Sprintf("tarball://root/repositories/%s.wiki.git", project.PathWithNamespace)
		repo.WikiURL = &wikiURL
	}
	r.exporter.State.Repositories = append(r.exporter.State.Repositories, repo)
}

func (r *RepositoryService) WriteFile() error {
	repos := r.exporter.State.Repositories
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Name > repos[j].Name
	})
	return r.exporter.WriteJsonFile(r.filename, repos)
}
