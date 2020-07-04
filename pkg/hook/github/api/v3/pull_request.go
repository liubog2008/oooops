package v3

import "time"

type PRLink struct {
	HRef string `json:"href,omitempty"`
}

type PRLinks struct {
	Self           *PRLink `json:"self,omitempty"`
	HTML           *PRLink `json:"html,omitempty"`
	Issue          *PRLink `json:"issue,omitempty"`
	Comments       *PRLink `json:"comments,omitempty"`
	ReviewComments *PRLink `json:"review_comments,omitempty"`
	ReviewComment  *PRLink `json:"review_comment,omitempty"`
	Commits        *PRLink `json:"commits,omitempty"`
	Statuses       *PRLink `json:"statuses,omitempty"`
}

type PullRequestBranch struct {
	Label string      `json:"label,omitempty"`
	Ref   string      `json:"ref,omitempty"`
	SHA   string      `json:"sha,omitempty"`
	Repo  *Repository `json:"repo,omitempty"`
	User  *User       `json:"user,omitempty"`
}

type PullRequest struct {
	ID                  int64      `json:"id,omitempty"`
	Number              int        `json:"number,omitempty"`
	State               string     `json:"state,omitempty"`
	Locked              bool       `json:"locked,omitempty"`
	Title               string     `json:"title,omitempty"`
	Body                string     `json:"body,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	ClosedAt            *time.Time `json:"closed_at,omitempty"`
	MergedAt            *time.Time `json:"merged_at,omitempty"`
	Labels              []Label    `json:"labels,omitempty"`
	User                *User      `json:"user,omitempty"`
	Draft               bool       `json:"draft,omitempty"`
	Merged              bool       `json:"merged,omitempty"`
	Mergeable           bool       `json:"mergeable,omitempty"`
	MergeableState      string     `json:"mergeable_state,omitempty"`
	MergedBy            *User      `json:"merged_by,omitempty"`
	MergeCommitSHA      string     `json:"merge_commit_sha,omitempty"`
	Rebaseable          bool       `json:"rebaseable,omitempty"`
	Comments            int        `json:"comments,omitempty"`
	Commits             int        `json:"commits,omitempty"`
	Additions           int        `json:"additions,omitempty"`
	Deletions           int        `json:"deletions,omitempty"`
	ChangedFiles        int        `json:"changed_files,omitempty"`
	URL                 string     `json:"url,omitempty"`
	HTMLURL             string     `json:"html_url,omitempty"`
	IssueURL            string     `json:"issue_url,omitempty"`
	StatusesURL         string     `json:"statuses_url,omitempty"`
	DiffURL             string     `json:"diff_url,omitempty"`
	PatchURL            string     `json:"patch_url,omitempty"`
	CommitsURL          string     `json:"commits_url,omitempty"`
	CommentsURL         string     `json:"comments_url,omitempty"`
	ReviewCommentsURL   string     `json:"review_comments_url,omitempty"`
	ReviewCommentURL    string     `json:"review_comment_url,omitempty"`
	ReviewComments      int        `json:"review_comments,omitempty"`
	Assignee            *User      `json:"assignee,omitempty"`
	Assignees           []User     `json:"assignees,omitempty"`
	Milestone           *Milestone `json:"milestone,omitempty"`
	MaintainerCanModify bool       `json:"maintainer_can_modify,omitempty"`
	AuthorAssociation   string     `json:"author_association,omitempty"`
	NodeID              string     `json:"node_id,omitempty"`
	RequestedReviewers  []User     `json:"requested_reviewers,omitempty"`

	RequestedTeams []Team `json:"requested_teams,omitempty"`

	Links *PRLinks           `json:"_links,omitempty"`
	Head  *PullRequestBranch `json:"head,omitempty"`
	Base  *PullRequestBranch `json:"base,omitempty"`

	ActiveLockReason string `json:"active_lock_reason,omitempty"`
}
