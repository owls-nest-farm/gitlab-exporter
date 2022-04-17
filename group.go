package main

import "github.com/xanzy/go-gitlab"

type GroupService struct {
	*BaseService
}

func NewGroupService(e *Exporter) *GroupService {
	return &GroupService{
		BaseService: &BaseService{
			exporter: e,
			filename: "groups.json",
		},
	}
}

func (g *GroupService) Get(gid string) (*gitlab.Group, *gitlab.Response, error) {
	return g.exporter.Client.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
}

func (g *GroupService) GetGroupProject(namespace, projectName string) ([]*gitlab.Project, *gitlab.Response, error) {
	f := false
	return g.exporter.Client.Groups.ListGroupProjects(namespace, &gitlab.ListGroupProjectsOptions{
		Archived: &f,
		Search:   &projectName,
	})
}

func (g *GroupService) GetSubgroups(gid interface{}) ([]*gitlab.Group, *gitlab.Response, error) {
	t := true
	return g.exporter.Client.Groups.ListSubGroups(gid, &gitlab.ListSubGroupsOptions{
		Owned: &t,
	})
}
