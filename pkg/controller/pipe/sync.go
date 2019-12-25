package pipe

import (
	"time"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func (c *Controller) syncPipe(key string) error {
	startTime := time.Now()

	defer func() {
		klog.V(4).Infof("Finished syncing pipe %q. (%v)", key, time.Since(startTime))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	pipe, err := c.pipeLister.Pipes(ns).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	events, err := c.ListWatchedEvents(pipe)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := c.GenerateFlow(pipe, &event); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) ListWatchedEvents(pipe *v1alpha1.Pipe) ([]v1alpha1.Event, error) {
	events, err := c.eventLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	watched := []v1alpha1.Event{}

	for _, e := range events {
		if isWatched(pipe, &e) {
			watched = append(watched, e)
		}
	}

	return watched, nil
}

func isWatched(pipe *v1alpha1.Pipe, event *v1alpha1.Event) bool {
	if event.Status.Phase == v1alpha1.EventConsumed {
		return false
	}
	if pipe.Spec.Git.Repo != event.Spec.Repo {
		return false
	}
	for _, on := range pipe.Spec.On {
		if on == event.Spec.When {
			return true
		}
	}
	return false
}

func (c *Controller) GenerateFlow(pipe *v1alpha1.Pipe, event *v1alpha1.Event) {
	klog.Infof("TODO: generate flow for pipe %s/%s, event: %s/%s", pipe.Namespace, pipe.Name, event.Namespac, event.Name)
	return nil
}
