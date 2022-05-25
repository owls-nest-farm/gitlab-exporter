package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

type ProjectService struct {
	*BaseService
}

type Project struct {
	Type      string     `json:"type"`
	URL       string     `json:"url"`
	Login     string     `json:"login"`
	Name      string     `json:"name"`
	Company   *string    `json:"company"`
	Website   string     `json:"website"`
	Location  *string    `json:"location"`
	Emails    []Email    `json:"emails"`
	CreatedAt *time.Time `json:"created_at"`
}

func NewProjectService(e *Exporter) *ProjectService {
	return &ProjectService{
		BaseService: &BaseService{
			exporter: e,
			filename: "projects.json",
		},
	}
}

//func (p *ProjectService) GetAll() ([]*gitlab.Project, *gitlab.Response, error) {
func (p *ProjectService) List() []int {
	return []int{
		35304986,
		35304985,
		35304984,
	}
	//	t := true
	//	f := false
	//	return p.exporter.Client.Projects.ListProjects(&gitlab.ListProjectsOptions{
	//		Archived:   &f, // Don't get deleted projects, which are considered archived.
	//		Membership: &t, // Only get projects to which the authenticated user has access.
	//	})
}

//func (p *ProjectService) ProjectIDs() []int {
//	projects, _, err := p.List()
//	if err != nil {
//		log.Fatalf("Failed to get user: %v", err)
//	}
//
//	ids := make([]int, len(projects))
//	for i, project := range projects {
//		ids[i] = project.ID
//	}
//	return ids
//
//	//	return Project{
//	//		Type:      "user",
//	//		URL:       user.WebURL,
//	//		Login:     user.Username,
//	//		Name:      user.Name,
//	//		Company:   nil,
//	//		Website:   user.WebsiteURL,
//	//		Location:  nil,
//	//		Emails:    getEmail(user.Email),
//	//		CreatedAt: user.CreatedAt,
//	//	}
//}

func (p *ProjectService) Get(pid int) (*gitlab.Project, *gitlab.Response, error) {
	return p.exporter.Client.Projects.GetProject(pid, &gitlab.GetProjectOptions{})
}

func (p *ProjectService) Export(pid int) {
	project, _, err := p.Get(pid)
	if err != nil {
		log.Fatalf("[ERROR] Failed to get project %d: %v", pid, err)
	}
	//	fmt.Println("project state", p.exporter.State)
	fmt.Println("project", project)
}
