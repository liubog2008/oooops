package v3

import "time"

type Label struct {
	ID          int64  `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Name        string `json:"name,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
}

type Milestone struct {
	URL          string     `json:"url,omitempty"`
	HTMLURL      string     `json:"html_url,omitempty"`
	LabelsURL    string     `json:"labels_url,omitempty"`
	ID           int64      `json:"id,omitempty"`
	Number       int        `json:"number,omitempty"`
	State        string     `json:"state,omitempty"`
	Title        string     `json:"title,omitempty"`
	Description  string     `json:"description,omitempty"`
	Creator      *User      `json:"creator,omitempty"`
	OpenIssues   int        `json:"open_issues,omitempty"`
	ClosedIssues int        `json:"closed_issues,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	DueOn        *time.Time `json:"due_on,omitempty"`
	NodeID       string     `json:"node_id,omitempty"`
}

type PullRequestLinks struct {
	URL      string `json:"url,omitempty"`
	HTMLURL  string `json:"html_url,omitempty"`
	DiffURL  string `json:"diff_url,omitempty"`
	PatchURL string `json:"patch_url,omitempty"`
}

type Reactions struct {
	TotalCount int    `json:"total_count,omitempty"`
	PlusOne    int    `json:"+1,omitempty"`
	MinusOne   int    `json:"-1,omitempty"`
	Laugh      int    `json:"laugh,omitempty"`
	Confused   int    `json:"confused,omitempty"`
	Heart      int    `json:"heart,omitempty"`
	Hooray     int    `json:"hooray,omitempty"`
	URL        string `json:"url,omitempty"`
}

type Issue struct {
	ID                int64             `json:"id,omitempty"`
	Number            int               `json:"number,omitempty"`
	State             string            `json:"state,omitempty"`
	Locked            bool              `json:"locked,omitempty"`
	Title             string            `json:"title,omitempty"`
	Body              string            `json:"body,omitempty"`
	AuthorAssociation string            `json:"author_association,omitempty"`
	User              *User             `json:"user,omitempty"`
	Labels            []Label           `json:"labels,omitempty"`
	Assignee          *User             `json:"assignee,omitempty"`
	Comments          int               `json:"comments,omitempty"`
	ClosedAt          time.Time         `json:"closed_at,omitempty"`
	CreatedAt         time.Time         `json:"created_at,omitempty"`
	UpdatedAt         time.Time         `json:"updated_at,omitempty"`
	ClosedBy          *User             `json:"closed_by,omitempty"`
	URL               string            `json:"url,omitempty"`
	HTMLURL           string            `json:"html_url,omitempty"`
	CommentsURL       string            `json:"comments_url,omitempty"`
	EventsURL         string            `json:"events_url,omitempty"`
	LabelsURL         string            `json:"labels_url,omitempty"`
	RepositoryURL     string            `json:"repository_url,omitempty"`
	Milestone         *Milestone        `json:"milestone,omitempty"`
	PullRequestLinks  *PullRequestLinks `json:"pull_request,omitempty"`
	Repository        *Repository       `json:"repository,omitempty"`
	Reactions         *Reactions        `json:"reactions,omitempty"`
	Assignees         []User            `json:"assignees,omitempty"`
	NodeID            string            `json:"node_id,omitempty"`

	// TextMatches is only populated from search results that request text matches
	// See: search.go and https://developer.github.com/v3/search/#text-match-metadata
	TextMatches []TextMatch `json:"text_matches,omitempty"`

	// ActiveLockReason is populated only when LockReason is provided while locking the issue.
	// Possible values are: "off-topic", "too heated", "resolved", and "spam".
	ActiveLockReason string `json:"active_lock_reason,omitempty"`
}
