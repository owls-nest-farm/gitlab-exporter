package main

import (
	"github.com/xanzy/go-gitlab"
)

type BranchService struct {
	*BaseService
}

type EmptyList []string

type Branch struct {
	Type                                 string    `json:"type"`
	Name                                 string    `json:"name"`
	URL                                  string    `json:"url"`
	CreatorURL                           string    `json:"creator_url"`
	RepositoryURL                        string    `json:"repository_url"`
	AdminEnforced                        bool      `json:"admin_enforced"`
	BlockDeletionsEnforcementLevel       int       `json:"block_deletions_enforcement_level"`
	BlockForcePushesEnforcementLevel     int       `json:"block_force_pushes_enforcement_level"`
	DismissStaleReviewsOnPush            bool      `json:"dismiss_stale_reviews_on_push"`
	PullRequestReviewsEnforcementLevel   string    `json:"pull_request_reviews_enforcement_level"`
	RequireCodeOwnerReview               bool      `json:"require_code_owner_review"`
	RequiredStatusChecksEnforcementLevel string    `json:"required_status_checks_enforcement_level"`
	StrictRequiredStatusChecksPolicy     bool      `json:"strict_required_status_checks_policy"`
	AuthorizedActorsOnly                 bool      `json:"authorized_actors_only"`
	AuthorizedUserUrls                   EmptyList `json:"authorized_user_urls"`
	AuthorizedTeamUrls                   EmptyList `json:"authorized_team_urls"`
	DismissalRestrictedUserUrls          EmptyList `json:"dismissal_restricted_user_urls"`
	DismissalRestrictedTeamUrls          EmptyList `json:"dismissal_restricted_team_urls"`
	RequiredStatusChecks                 EmptyList `json:"required_status_checks"`
}

func NewBranchService(e *Exporter) *BranchService {
	return &BranchService{
		BaseService: &BaseService{
			exporter: e,
			filename: "protected_branches.json",
		},
	}
}
func (b *BranchService) GetAll() []Branch {
	project := b.exporter.CurrentProject
	branches, _, err := b.exporter.Client.ProtectedBranches.ListProtectedBranches(project.ID, &gitlab.ListProtectedBranchesOptions{})
	if err != nil {
		panic(err)
	}

	br := make([]Branch, len(branches))
	for i, branch := range branches {
		br[i] = Branch{
			Type:                                 "protected_branch",
			Name:                                 branch.Name,
			URL:                                  project.WebURL,
			CreatorURL:                           b.exporter.CurrentUser.URL,
			RepositoryURL:                        project.WebURL,
			AdminEnforced:                        true,
			BlockDeletionsEnforcementLevel:       2,
			BlockForcePushesEnforcementLevel:     2,
			DismissStaleReviewsOnPush:            false,
			PullRequestReviewsEnforcementLevel:   "off",
			RequiredStatusChecksEnforcementLevel: "off",
			StrictRequiredStatusChecksPolicy:     false,
			AuthorizedActorsOnly:                 false,
			AuthorizedUserUrls:                   EmptyList{},
			AuthorizedTeamUrls:                   EmptyList{},
			DismissalRestrictedUserUrls:          EmptyList{},
			DismissalRestrictedTeamUrls:          EmptyList{},
			RequiredStatusChecks:                 EmptyList{},
		}
	}

	return br
}

func (b *BranchService) Export() {
	//	cc, err := c.GetAll()
	//	fmt.Println("cc", cc, err)
	//	if err != nil {
	//		log.Fatalf("Failed to get commit comments for projectID %d: %v", cc.exporter.CurrentProject.ID, err)
	//	}
	branches := b.GetAll()
	if len(branches) > 0 {
		b.exporter.State.Branches = append(b.exporter.State.Branches, branches...)
	}
}

func (b *BranchService) WriteFile() error {
	return b.exporter.WriteJsonFile(b.filename, b.exporter.State.Branches)
}
