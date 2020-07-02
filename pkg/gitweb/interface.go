package gitweb

import (
	"net/http"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

// Interface defines action to handle webhook of different git
// platform like github, gitlab
type Interface interface {
	ParseWebhook(r *http.Request) (*v1alpha1.Event, error)
}
