package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/xanzy/go-gitlab"
)

// CommitComment represents a GitLab commit comment.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type CommitComment struct {
	Note      string     `json:"note"`
	Path      string     `json:"path"`
	Line      int        `json:"line"`
	LineType  string     `json:"line_type"`
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

func getCommits(project *gitlab.Project) ([]*gitlab.Commit, *gitlab.Response, error) {
	return getClient().Commits.ListCommits(project.ID, &gitlab.ListCommitsOptions{})
}

func getCommitComments(project *gitlab.Project) ([]CommitComment, error) {
	// The `xanzy/go-gitlab` SDK has structs that don't match the response, so we're not getting
	// several fields that we are essential (`CommitComment.CreatedAt`, `Author.WebURL`).
	// So, we need to fetch this ourselves with the correct structs :(
	api := fmt.Sprintf("projects/%d/repository/commits/%s/comments", project.ID, url.PathEscape(project.DefaultBranch))
	resp, err := newRequest(api)
	if err != nil {
		return nil, err
	}

	var cc []CommitComment
	err = json.NewDecoder(resp.Body).Decode(&cc)

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

	return cc, err
}
