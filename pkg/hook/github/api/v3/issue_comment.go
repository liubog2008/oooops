package v3

import "time"

type IssueComment struct {
	ID                int64      `json:"id,omitempty"`
	NodeID            string     `json:"node_id,omitempty"`
	Body              string     `json:"body,omitempty"`
	User              *User      `json:"user,omitempty"`
	Reactions         *Reactions `json:"reactions,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	AuthorAssociation string     `json:"author_association,omitempty"`
	URL               string     `json:"url,omitempty"`
	HTMLURL           string     `json:"html_url,omitempty"`
	IssueURL          string     `json:"issue_url,omitempty"`
}
