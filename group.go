package main

import (
	"github.com/xanzy/go-gitlab"
)

type GroupService struct {
	*BaseService
}

func NewGroupService(e *Exporter) *GroupService {
	return &GroupService{
		BaseService: &BaseService{
			exporter: e,
			filename: "organizations.json",
		},
	}
}

type Group struct {
	Type        string   `json:"type"`
	URL         string   `json:"url"`
	Login       string   `json:"login"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Website     *string  `json:"website"`
	Location    *string  `json:"location"`
	Email       *string  `json:"email"`
	Members     []string `json:"members"`
}

func (g *GroupService) Get(gid string) Group {
	group, _, err := g.exporter.Client.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
	if err != nil {
		panic(err)
	}

	//	fmt.Println(g.exporter.CurrentProject.Namespace)
	return Group{
		Type:        "organization",
		URL:         group.WebURL,
		Login:       group.Path,
		Name:        group.Name,
		Description: group.Description,
	}
}

func (g *GroupService) GetGroupProject(namespace, projectName string) ([]*gitlab.Project, *gitlab.Response, error) {
	f := false
	return g.exporter.Client.Groups.ListGroupProjects(namespace, &gitlab.ListGroupProjectsOptions{
		Archived: &f,
		Search:   &projectName,
	})
}

//func (g *GroupService) GetSubgroups(gid interface{}) ([]*gitlab.Group, *gitlab.Response, error) {
//	t := true
//	return g.exporter.Client.Groups.ListSubGroups(gid, &gitlab.ListSubGroupsOptions{
//		Owned: &t,
//	})
//}

func (g *GroupService) Export(gid string) {
	//	cc, err := c.GetAll()
	//	fmt.Println("cc", cc, err)
	//	if err != nil {
	//		log.Fatalf("Failed to get commit comments for projectID %d: %v", cc.exporter.CurrentProject.ID, err)
	//	}
	group := g.Get(gid)

	var contains bool
	for _, g := range g.exporter.State.Groups {
		if g.Login == group.Login {
			contains = true
			break
		}
	}

	if !contains {
		g.exporter.State.Groups = append(g.exporter.State.Groups, group)
	}
}

func (g *GroupService) WriteFile() error {
	return g.exporter.WriteJsonFile(g.filename, g.exporter.State.Groups)
}
