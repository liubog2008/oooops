package v3

import "time"

type InstallationPermissions struct {
	Administration              string `json:"administration,omitempty"`
	Blocking                    string `json:"blocking,omitempty"`
	Checks                      string `json:"checks,omitempty"`
	Contents                    string `json:"contents,omitempty"`
	ContentReferences           string `json:"content_references,omitempty"`
	Deployments                 string `json:"deployments,omitempty"`
	Emails                      string `json:"emails,omitempty"`
	Followers                   string `json:"followers,omitempty"`
	Issues                      string `json:"issues,omitempty"`
	Metadata                    string `json:"metadata,omitempty"`
	Members                     string `json:"members,omitempty"`
	OrganizationAdministration  string `json:"organization_administration,omitempty"`
	OrganizationHooks           string `json:"organization_hooks,omitempty"`
	OrganizationPlan            string `json:"organization_plan,omitempty"`
	OrganizationPreReceiveHooks string `json:"organization_pre_receive_hooks,omitempty"`
	OrganizationProjects        string `json:"organization_projects,omitempty"`
	OrganizationUserBlocking    string `json:"organization_user_blocking,omitempty"`
	Packages                    string `json:"packages,omitempty"`
	Pages                       string `json:"pages,omitempty"`
	PullRequests                string `json:"pull_requests,omitempty"`
	RepositoryHooks             string `json:"repository_hooks,omitempty"`
	RepositoryProjects          string `json:"repository_projects,omitempty"`
	RepositoryPreReceiveHooks   string `json:"repository_pre_receive_hooks,omitempty"`
	SingleFile                  string `json:"single_file,omitempty"`
	Statuses                    string `json:"statuses,omitempty"`
	TeamDiscussions             string `json:"team_discussions,omitempty"`
	VulnerabilityAlerts         string `json:"vulnerability_alerts,omitempty"`
}

// Installation represents a GitHub Apps installation.
type Installation struct {
	ID                  int64                    `json:"id,omitempty"`
	AppID               int64                    `json:"app_id,omitempty"`
	TargetID            int64                    `json:"target_id,omitempty"`
	Account             *User                    `json:"account,omitempty"`
	AccessTokensURL     string                   `json:"access_tokens_url,omitempty"`
	RepositoriesURL     string                   `json:"repositories_url,omitempty"`
	HTMLURL             string                   `json:"html_url,omitempty"`
	TargetType          string                   `json:"target_type,omitempty"`
	SingleFileName      string                   `json:"single_file_name,omitempty"`
	RepositorySelection string                   `json:"repository_selection,omitempty"`
	Events              []string                 `json:"events,omitempty"`
	Permissions         *InstallationPermissions `json:"permissions,omitempty"`
	CreatedAt           time.Time                `json:"created_at,omitempty"`
	UpdatedAt           time.Time                `json:"updated_at,omitempty"`
}
