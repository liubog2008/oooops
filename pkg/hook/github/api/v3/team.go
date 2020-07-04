package v3

type Team struct {
	ID          int64  `json:"id,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Slug        string `json:"slug,omitempty"`

	Permission string `json:"permission,omitempty"`

	Privacy string `json:"privacy,omitempty"`

	MembersCount    int           `json:"members_count,omitempty"`
	ReposCount      int           `json:"repos_count,omitempty"`
	Organization    *Organization `json:"organization,omitempty"`
	MembersURL      string        `json:"members_url,omitempty"`
	RepositoriesURL string        `json:"repositories_url,omitempty"`
	Parent          *Team         `json:"parent,omitempty"`

	LDAPDN string `json:"ldap_dn,omitempty"`
}
