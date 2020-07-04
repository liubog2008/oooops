package v3

import "time"

// Plan represents the payment plan for an account. See plans at https://github.com/plans.
type Plan struct {
	Name          string `json:"name,omitempty"`
	Space         int    `json:"space,omitempty"`
	Collaborators int    `json:"collaborators,omitempty"`
	PrivateRepos  int    `json:"private_repos,omitempty"`
	FilledSeats   int    `json:"filled_seats,omitempty"`
	Seats         int    `json:"seats,omitempty"`
}

// Match represents a single text match.
type Match struct {
	Text    string `json:"text,omitempty"`
	Indices []int  `json:"indices,omitempty"`
}

// TextMatch represents a text match for a SearchResult
type TextMatch struct {
	ObjectURL  string  `json:"object_url,omitempty"`
	ObjectType string  `json:"object_type,omitempty"`
	Property   string  `json:"property,omitempty"`
	Fragment   string  `json:"fragment,omitempty"`
	Matches    []Match `json:"matches,omitempty"`
}

// User represents a GitHub user.
type User struct {
	Login                   string    `json:"login,omitempty"`
	ID                      int64     `json:"id,omitempty"`
	NodeID                  string    `json:"node_id,omitempty"`
	AvatarURL               string    `json:"avatar_url,omitempty"`
	HTMLURL                 string    `json:"html_url,omitempty"`
	GravatarID              string    `json:"gravatar_id,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Company                 string    `json:"company,omitempty"`
	Blog                    string    `json:"blog,omitempty"`
	Location                string    `json:"location,omitempty"`
	Email                   string    `json:"email,omitempty"`
	Hireable                bool      `json:"hireable,omitempty"`
	Bio                     string    `json:"bio,omitempty"`
	TwitterUsername         string    `json:"twitter_username,omitempty"`
	PublicRepos             int       `json:"public_repos,omitempty"`
	PublicGists             int       `json:"public_gists,omitempty"`
	Followers               int       `json:"followers,omitempty"`
	Following               int       `json:"following,omitempty"`
	CreatedAt               time.Time `json:"created_at,omitempty"`
	UpdatedAt               time.Time `json:"updated_at,omitempty"`
	SuspendedAt             time.Time `json:"suspended_at,omitempty"`
	Type                    string    `json:"type,omitempty"`
	SiteAdmin               bool      `json:"site_admin,omitempty"`
	TotalPrivateRepos       int       `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos       int       `json:"owned_private_repos,omitempty"`
	PrivateGists            int       `json:"private_gists,omitempty"`
	DiskUsage               int       `json:"disk_usage,omitempty"`
	Collaborators           int       `json:"collaborators,omitempty"`
	TwoFactorAuthentication bool      `json:"two_factor_authentication,omitempty"`
	Plan                    *Plan     `json:"plan,omitempty"`
	LdapDn                  string    `json:"ldap_dn,omitempty"`

	// API URLs
	URL               string `json:"url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`

	// TextMatches is only populated from search results that request text matches
	// See: search.go and https://developer.github.com/v3/search/#text-match-metadata
	TextMatches []TextMatch `json:"text_matches,omitempty"`

	// Permissions identifies the permissions that a user has on a given
	// repository. This is only populated when calling Repositories.ListCollaborators.
	Permissions map[string]bool `json:"permissions,omitempty"`
}
