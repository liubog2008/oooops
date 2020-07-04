package github

import (
	"k8s.io/klog"

	v3 "github.com/liubog2008/oooops/pkg/hook/github/api/v3"
)

func (h *githubHook) handleIssueCommentEvent(e *v3.IssueCommentEvent) {
	defer h.wg.Done()
	klog.Infof("handle issue comment event: \n\n%v\n", e)
}

func (h *githubHook) handlePullRequestEvent(e *v3.PullRequestEvent) {
	defer h.wg.Done()
	klog.Infof("handle pull request event: \n\n%v\n", e)
}

func (h *githubHook) handleCreateEvent(e *v3.CreateEvent) {
	defer h.wg.Done()
	klog.Infof("handle create event: \n\n%v\n", e)
}
