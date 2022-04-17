package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"
)

type CommitCommentService struct {
	*BaseService
}

// CommitComment represents a GitLab commit comment.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type CommitComment struct {
	Type      string     `json:"type"`
	Note      string     `json:"note"`
	Path      *string    `json:"path"`
	Line      *int       `json:"line"`
	LineType  *string    `json:"line_type"`
	Author    Author     `json:"author"`
	CreatedAt *time.Time `json:"created_at"`
}

// Author represents a GitLab commit author
type Author struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func NewCommitCommentService(e *Exporter) *CommitCommentService {
	return &CommitCommentService{
		BaseService: &BaseService{
			exporter: e,
			filename: "commit_comments.json",
		},
	}
}

//[{"note":"this was the original commit!","path":null,"line":null,"line_type":null,"author":{"id":10924637,"username":"btoll","name":"btoll","state":"active","avatar_url":"https://secure.gravatar.com/avatar/3a614474bb70fa358ba728a058ed9768?s=80\u0026d=identicon","web_url":"https://gitlab.com/btoll"},"created_at":"2022-04-16T18:15:02.206Z"}]

//type CommitComment struct {
//	Type       string     `json:"type"`
//	URL        string     `json:"url"`
//	Repository string     `json:"repository"`
//	User       string     `json:"user"`
//	Body       string     `json:"body"`
//	Formatter  string     `json:"formatter"`
//	Path       *string    `json:"path"`
//	Position   *int       `json:"position"`
//	CommitID   string     `json:"commit_id"`
//	CreatedAt  *time.Time `json:"created_at"`
//}

// [{"note":"this was the original commit!","path":null,"line":null,"line_type":null,"author":{"id":10924637,"username":"btoll","name":"btoll","state":"active","avatar_url":"https://secure.gravatar.com/avatar/3a614474bb70fa358ba728a058ed9768?s=80\u0026d=identicon","web_url":"https://gitlab.com/btoll"},"created_at":"2022-04-16T18:15:02.206Z"}]

//$ curl -X GET --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/projects/35304986/repository/commits/master/comments
//[{"note":"this was the original commit!","path":null,"line":null,"line_type":null,"author":{"id":10924637,"username":"btoll","name":"btoll","state":"active","avatar_url":"https://secure.gravatar.com/avatar/3a614474bb70fa358ba728a058ed9768?s=80\u0026d=identicon","web_url":"https://gitlab.com/btoll"},"created_at":"2022-04-16T18:15:02.206Z"}]

//func (ccs *CommitCommentService) Get(project *gitlab.Project) ([]*gitlab.Commit, *gitlab.Response, error) {
//	return ccs.exporter.Client.Commits.ListCommits(project.ID, &gitlab.ListCommitsOptions{})
//}

func (c *CommitCommentService) GetAll() []CommitComment {
	project := c.exporter.CurrentProject

	// The `xanzy/go-gitlab` SDK has structs that don't match the response, so we're not getting
	// several fields that we are essential (`CommitComment.CreatedAt`, `Author.WebURL`).
	// So, we need to fetch this ourselves with the correct structs :(
	api := fmt.Sprintf("projects/%d/repository/commits/%s/comments", project.ID, url.PathEscape(project.DefaultBranch))
	resp, err := c.exporter.NewRequest(api)
	if err != nil {
		log.Fatalf("Failed to create the request object for the commit comments for projectID %d: %v", project.ID, err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the commit comments for projectID %d: %v", project.ID, err)
	}
	defer resp.Body.Close()

	var ccs []CommitComment
	// If there's an error unmarshalling the response body into the structs, it's ok (for now),
	// it will return an empty slice.
	json.Unmarshal(data, &ccs)
	return ccs

	//	err = json.NewDecoder(resp.Body).Decode(&c)

	//	git := getClient()
	//	commitComments, _, err := git.Commits.GetCommitComments(project.ID, project.DefaultBranch, &gitlab.GetCommitCommentsOptions{})
	//	if err != nil {
	//		log.Fatalf("Failed to get commit comments for projectID %d: %v", project.ID, err)
	//	}

	//	for i, commitComment := range cc {
	//		cc[i] = CommitComment{
	//			Type:       "commit_comment",
	//			URL:        project.WebURL,
	//			Repository: project.WebURL,
	//			User:       commitComment.Author.Username,
	//			Body:       commitComment.Note,
	//			Formatter:  "markdown",
	//			Path:       &commitComment.Path,
	//			Position:   &commitComment.Line,
	//			//			CommitID:   commitComment.ID,
	//			CreatedAt: commitComment.Author.CreatedAt,
	//		}
	//	}
}

func (c *CommitCommentService) Export() {
	//	cc, err := c.GetAll()
	//	fmt.Println("cc", cc, err)
	//	if err != nil {
	//		log.Fatalf("Failed to get commit comments for projectID %d: %v", cc.exporter.CurrentProject.ID, err)
	//	}
	comments := c.GetAll()
	if len(comments) > 0 {
		c.exporter.State.CommitComments = append(c.exporter.State.CommitComments, comments...)
	}
}

func (c *CommitCommentService) WriteFile() error {
	return c.exporter.WriteJsonFile(c.filename, c.exporter.State.CommitComments)
}
