package v3

type EventType string

const (
	EventTypePullRequest  EventType = "pull_request"
	EventTypeIssueComment EventType = "issue_comment"
	EventTypeCreate       EventType = "create"
)

const (
	// HeaderGithubEvent defines name of the event that triggered the delivery
	HeaderGithubEvent = "X-Github-Event"
	// HeaderGithubDelivery defines a GUID to identify the delivery
	HeaderGithubDelivery = "X-Github-Delivery"
	// HeaderHubSignature is the HMAC hex digest of the response body
	// This header will be sent if the webhook is configured with a secret
	// The HMAC hex digest is generated using the sha1 hash function and
	// the secret as the HMAC key
	HeaderHubSignature = "X-Hub-Signature"

	// HeaderUserAgentPrefix defines prefix of webhook request User-Agent
	HeaderUserAgentPrefix = "GitHub-Hookshot/"
)

type Common struct {
	// GUID from X-Github-Delivery
	GUID string

	Sender       *User         `json:"sender,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

type EditChangeFrom struct {
	From string `json:"from,omitempty"`
}

type EditChange struct {
	Title *EditChangeFrom `json:"title,omitempty"`
	Body  *EditChangeFrom `json:"body,omitempty"`
}

type CheckRunAction string

const (
	CheckRunCreated         CheckRunAction = "created"
	CheckRunCompleted       CheckRunAction = "completed"
	CheckRunRerequested     CheckRunAction = "rerequested"
	CheckRunRequestedAction CheckRunAction = "requested_action"
)

type RequestedAction struct {
	Identifier string `json:"identifier"`
}

// CheckRunEvent represents a check run activity has occurred
// The type of activity is specified in the action property of
// the payload object
// For more information, see the "check runs" REST API.
//
// Github API docs: https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads#check_run
// type CheckRunEvent struct {
// 	Common `json:",inline"`
//
// 	Action CheckRunAction `json:"action"`
//
// 	CheckRun        *CheckRun        `json:"check_run,omitempty"`
// 	RequestedAction *RequestedAction `json:"requested_action,omitempty"`
// }

type CreateRefType string

const (
	CreateRefTypeTag    CreateRefType = "tag"
	CreateRefTypeBranch CreateRefType = "branch"
)

// CreateEvent represents a git branch or tag is created
// For more information, see the "Git data" REST API
//
// Github API docs: https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads#create
type CreateEvent struct {
	Common       `json:",inline"`
	Ref          string        `json:"ref,omitempty"`
	RefType      CreateRefType `json:"ref_type,omitempty"`
	MasterBranch string        `json:"master_branch,omitempty"`
	Description  string        `json:"description,omitempty"`
}

type IssueCommentAction string

const (
	IssueCommentCreated IssueCommentAction = "created"
	IssueCommentEdited  IssueCommentAction = "edited"
	IssueCommentDeleted IssueCommentAction = "deleted"
)

// IssueCommentEvent is triggered by activity related to an issue comment
// The type of activity is specified in the action property of
// the payload object
// For more information, see the "issue comments" REST API.
//
// GitHub API docs: https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads#issue_comment
type IssueCommentEvent struct {
	Common `json:",inline"`
	// Action is the action that was performed on the comment.
	// Possible values are: "created", "edited", "deleted".
	Action  IssueCommentAction `json:"action,omitempty"`
	Issue   *Issue             `json:"issue,omitempty"`
	Comment *IssueComment      `json:"comment,omitempty"`

	Changes *EditChange `json:"changes,omitempty"`
}

type PullRequestAction string

const (
	PullRequestOpened               PullRequestAction = "opened"
	PullRequestEdited               PullRequestAction = "edited"
	PullRequestClosed               PullRequestAction = "closed"
	PullRequestAssigned             PullRequestAction = "assigned"
	PullRequestUnassigned           PullRequestAction = "unassigned"
	PullRequestReviewRequested      PullRequestAction = "review_requested"
	PullRequestReviewRequestRemoved PullRequestAction = "review_request_removed"
	PullRequestReadyForReview       PullRequestAction = "ready_for_review"
	PullRequestLabeled              PullRequestAction = "labeled"
	PullRequestUnlabeled            PullRequestAction = "unlabeled"
	PullRequestSynchronize          PullRequestAction = "synchronize"
	PullRequestLocked               PullRequestAction = "locked"
	PullRequestUnlocked             PullRequestAction = "unlocked"
	PullRequestReopened             PullRequestAction = "reopened"
)

// PullRequestEvent is triggered by activity related to pull requests
// The type of activity is specified in the action property of
// the payload object
// For more information, see the "pull requests" REST API

// Github API docs: https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads#pull_request
type PullRequestEvent struct {
	Common `json:",inline"`

	Action      PullRequestAction `json:"action"`
	Number      int               `json:"number"`
	PullRequest *PullRequest      `json:"pull_request,omitempty"`

	Changes *EditChange `json:"changes,omitempty"`
}
