package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

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

func listProjects() ([]*gitlab.Project, *gitlab.Response, error) {
	t := true
	f := false
	return getClient().Projects.ListProjects(&gitlab.ListProjectsOptions{
		Archived:   &f, // Don't get deleted projects, which are considered archived.
		Membership: &t, // Only get projects to which the authenticated user has access.
	})
}

func getProjectIDs() []int {
	projects, _, err := listProjects()
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	ids := make([]int, len(projects))
	for i, project := range projects {
		ids[i] = project.ID
	}
	return ids

	//	return Project{
	//		Type:      "user",
	//		URL:       user.WebURL,
	//		Login:     user.Username,
	//		Name:      user.Name,
	//		Company:   nil,
	//		Website:   user.WebsiteURL,
	//		Location:  nil,
	//		Emails:    getEmail(user.Email),
	//		CreatedAt: user.CreatedAt,
	//	}
}

func getProjects() {
	git := getClient()
	for _, projectID := range getProjectIDs() {
		project, _, err := git.Projects.GetProject(projectID, &gitlab.GetProjectOptions{})
		if err != nil {
			log.Fatalf("Failed to get projectID %d: %v", projectID, err)
		}

		//		getRepository(project)
		//		labels, _, err := getLabels(project)
		//		if err != nil {
		//			log.Fatalf("Failed to get labels for projectID %d: %v", projectID, err)
		//		}
		//		fmt.Println(labels)

		//		ps := getIssues(project)
		//		if len(ps) > 0 {
		//			file, _ := json.Marshal(ps)
		//			ioutil.WriteFile("issues.json", file, 0644)
		//		}
		//
		//		prs := getPullRequests(project)
		//		if len(prs) > 0 {
		//			file, _ := json.Marshal(prs)
		//			ioutil.WriteFile("pull_requests.json", file, 0644)
		//		}
		//		commits, _, err := getCommits(project)
		//
		//		if err != nil {
		//			log.Fatalf("Failed to get commits for projectID %d: %v", project.ID, err)
		//		}
		//		res, _ := getCommitComments(project)
		//		fmt.Println("res", res)
		res2, _ := getWebhooks(project)
		fmt.Println("res2", res2)
		//		res3, _ := getTags(project)
		//		fmt.Println("res3", res3)
		//		ccs := getCommitComments(project, commits)
		//		if project.ID == 35304986 {
		//			if len(ccs) > 0 {
		//				fmt.Printf("%+v\n", ccs)
		//				file, _ := json.Marshal(ccs)
		//				ioutil.WriteFile("commit_comments.json", file, 0644)
		//			}
		//		}
		//		fmt.Println("commits", len(commits))
		//		fmt.Println("projectID", project.ID)
	}
}
