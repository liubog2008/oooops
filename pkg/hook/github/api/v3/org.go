package v3

import "time"

type Organization struct {
	Login                       string    `json:"login,omitempty"`
	ID                          int64     `json:"id,omitempty"`
	NodeID                      string    `json:"node_id,omitempty"`
	AvatarURL                   string    `json:"avatar_url,omitempty"`
	HTMLURL                     string    `json:"html_url,omitempty"`
	Name                        string    `json:"name,omitempty"`
	Company                     string    `json:"company,omitempty"`
	Blog                        string    `json:"blog,omitempty"`
	Location                    string    `json:"location,omitempty"`
	Email                       string    `json:"email,omitempty"`
	TwitterUsername             string    `json:"twitter_username,omitempty"`
	Description                 string    `json:"description,omitempty"`
	PublicRepos                 int       `json:"public_repos,omitempty"`
	PublicGists                 int       `json:"public_gists,omitempty"`
	Followers                   int       `json:"followers,omitempty"`
	Following                   int       `json:"following,omitempty"`
	CreatedAt                   time.Time `json:"created_at,omitempty"`
	UpdatedAt                   time.Time `json:"updated_at,omitempty"`
	TotalPrivateRepos           int       `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos           int       `json:"owned_private_repos,omitempty"`
	PrivateGists                int       `json:"private_gists,omitempty"`
	DiskUsage                   int       `json:"disk_usage,omitempty"`
	Collaborators               int       `json:"collaborators,omitempty"`
	BillingEmail                string    `json:"billing_email,omitempty"`
	Type                        string    `json:"type,omitempty"`
	Plan                        *Plan     `json:"plan,omitempty"`
	TwoFactorRequirementEnabled bool      `json:"two_factor_requirement_enabled,omitempty"`
	IsVerified                  bool      `json:"is_verified,omitempty"`
	HasOrganizationProjects     bool      `json:"has_organization_projects,omitempty"`
	HasRepositoryProjects       bool      `json:"has_repository_projects,omitempty"`

	DefaultRepoPermission string `json:"default_repository_permission,omitempty"`
	DefaultRepoSettings   string `json:"default_repository_settings,omitempty"`

	MembersCanCreateRepos bool `json:"members_can_create_repositories,omitempty"`

	MembersCanCreatePublicRepos   bool `json:"members_can_create_public_repositories,omitempty"`
	MembersCanCreatePrivateRepos  bool `json:"members_can_create_private_repositories,omitempty"`
	MembersCanCreateInternalRepos bool `json:"members_can_create_internal_repositories,omitempty"`

	MembersAllowedRepositoryCreationType string `json:"members_allowed_repository_creation_type,omitempty"`

	URL              string `json:"url,omitempty"`
	EventsURL        string `json:"events_url,omitempty"`
	HooksURL         string `json:"hooks_url,omitempty"`
	IssuesURL        string `json:"issues_url,omitempty"`
	MembersURL       string `json:"members_url,omitempty"`
	PublicMembersURL string `json:"public_members_url,omitempty"`
	ReposURL         string `json:"repos_url,omitempty"`
}
