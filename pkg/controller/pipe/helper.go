package pipe

import (
	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

func isWatched(pipe *v1alpha1.Pipe, event *v1alpha1.Event) bool {
	if pipe.Spec.Git.Repo != event.Spec.Repo {
		return false
	}
	for _, when := range pipe.Spec.When {
		if when == event.Spec.When {
			return true
		}
	}
	return false
}

func (c *Controller) getPipeWatchers(event *v1alpha1.Event) ([]*v1alpha1.Pipe, error) {
	pipes, err := c.pipeLister.Pipes(event.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	watchers := []*v1alpha1.Pipe{}
	for _, watcher := range pipes {
		if isWatched(watcher, event) {
			watchers = append(watchers, watcher)
		}
	}
	return watchers, nil

}
