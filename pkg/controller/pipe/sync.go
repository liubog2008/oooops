package pipe

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/utils/random"
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
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	events, err := c.listWatchedEvents(pipe)
	if err != nil {
		return err
	}

	for _, event := range events {
		klog.V(6).Infof("consume event: %v/%v", event.Namespace, event.Name)
		if err := c.generateFlow(pipe, event); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) listWatchedEvents(pipe *v1alpha1.Pipe) ([]*v1alpha1.Event, error) {
	events, err := c.eventLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	watched := []*v1alpha1.Event{}

	for _, e := range events {
		if isWatched(pipe, e) {
			watched = append(watched, e)
		}
	}

	klog.V(6).Infof("event total: %v, watched by pipe %v: %v", len(events), pipe.Name, len(watched))

	return watched, nil
}

func (c *Controller) generateFlow(pipe *v1alpha1.Pipe, event *v1alpha1.Event) error {
	klog.Infof("TODO: generate flow for pipe %s/%s, event: %s/%s", pipe.Namespace, pipe.Name, event.Namespace, event.Name)

	ns := pipe.Namespace
	hashCode := hash(event)

	flows, err := c.flowLister.Flows(ns).List(labels.SelectorFromValidatedSet(
		labels.Set{
			v1alpha1.DefaultFlowRevisionLabelKey: hashCode,
		},
	))
	if err != nil {
		return err
	}
	name := genName(pipe)

	pipeSpec := pipe.Spec.DeepCopy()
	selector := pipeSpec.Selector.DeepCopy()

	if selector == nil || selector.MatchLabels == nil {
		selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{},
		}
	}
	selector.MatchLabels[v1alpha1.DefaultFlowRevisionLabelKey] = hashCode

	owner := metav1.NewControllerRef(pipe, c.GroupVersionKind)

	expectedFlow := &v1alpha1.Flow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				v1alpha1.DefaultFlowRevisionLabelKey: hashCode,
			},
			OwnerReferences: []metav1.OwnerReference{
				*owner,
			},
		},
		Spec: v1alpha1.FlowSpec{
			Selector: selector,
			Git:      pipeSpec.Git,
			Stages:   pipeSpec.Stages,
		},
		Status: v1alpha1.FlowStatus{
			Phase: v1alpha1.FlowPending,
		},
	}
	expectedFlow.Spec.Git.Ref = event.Spec.Ref

	// no flow is selected
	if len(flows) == 0 {
		if _, err := c.extClient.MarioV1alpha1().Flows(ns).Create(expectedFlow); err != nil {
			return err
		}
		return nil
	}

	for _, flow := range flows {
		// ignore flow which is not controlled by this pipe
		if !metav1.IsControlledBy(flow, pipe) {
			return nil
		}
		// ignore flow which is not triggered by this event
		if !isTriggeredBy(flow, event) {
			return nil
		}

		if !sementicEqual(flow, expectedFlow) {
			updating := flow.DeepCopy()
			updating.Spec.Mario = nil
			updating.Spec.Git = expectedFlow.Spec.Git
			updating.Spec.Stages = expectedFlow.Spec.Stages

			updating.Status.Phase = v1alpha1.FlowPending

			if _, err := c.extClient.MarioV1alpha1().Flows(ns).Update(updating); err != nil {
				return nil
			}
		}
	}

	return nil
}

func sementicEqual(a, b *v1alpha1.Flow) bool {
	if !reflect.DeepEqual(&a.Spec.Git, &b.Spec.Git) {
		return false
	}
	if !reflect.DeepEqual(&a.Spec.Stages, &b.Spec.Stages) {
		return false
	}
	return true
}

// hash generate event identity
// NOTE(liubog2008): maybe change to event UID?
func hash(event *v1alpha1.Event) string {
	repo := event.Spec.Repo
	ref := event.Spec.Ref

	hasher := md5.New()
	hasher.Write([]byte(repo + "@" + ref))
	code := hex.EncodeToString(hasher.Sum(nil))
	return code[:11]
}

// genName generates name of flow
func genName(pipe *v1alpha1.Pipe) string {
	return pipe.Name + "-" + random.Random(7)
}

func isTriggeredBy(flow *v1alpha1.Flow, event *v1alpha1.Event) bool {
	if flow.Spec.Git.Repo != event.Spec.Repo {
		return false
	}

	if flow.Spec.Git.Ref != event.Spec.Ref {
		return false
	}

	return true
}
