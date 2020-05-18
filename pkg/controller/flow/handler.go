package flow

import (
	"fmt"
	"reflect"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) addJob(obj interface{}) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("obj is not Job: %v", obj))
		return
	}
	ref := metav1.GetControllerOf(job)
	if ref == nil {
		return
	}

	flow := c.getFlowFromRef(job.Namespace, ref)
	if flow == nil {
		return
	}

	c.addFlow(flow)
}

func (c *Controller) updateJob(old, cur interface{}) {
	oldJob, ok1 := old.(*batchv1.Job)
	curJob, ok2 := cur.(*batchv1.Job)
	if !ok1 || !ok2 {
		utilruntime.HandleError(fmt.Errorf("either old or cur is not Job: %v, %v", old, cur))
		return
	}
	if oldJob.ResourceVersion == curJob.ResourceVersion {
		return
	}

	oldRef := metav1.GetControllerOf(oldJob)
	curRef := metav1.GetControllerOf(curJob)

	refChanged := reflect.DeepEqual(oldRef, curRef)

	if oldRef != nil && refChanged {
		if oldFlow := c.getFlowFromRef(oldJob.Namespace, oldRef); oldFlow != nil {
			c.addFlow(oldFlow)
		}
	}
	if curRef != nil {
		if curFlow := c.getFlowFromRef(curJob.Namespace, curRef); curFlow != nil {
			c.addFlow(curFlow)
		}
	}
	// When job's owner is deleted, flow status should be updated
	// TODO(liubog2008): add logic to deal ref: not nil -> nil
}

func (c *Controller) deleteJob(obj interface{}) {
	if job, ok := obj.(*batchv1.Job); ok {
		c.addJob(job)
		return
	}
	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
		return
	}
	job, ok := tombstone.Obj.(*batchv1.Job)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not an Event: %#v", tombstone.Obj))
		return
	}
	c.addJob(job)
}

func (c *Controller) addFlow(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}

	c.queue.Add(key)
}

func (c *Controller) updateFlow(old, cur interface{}) {
	c.addFlow(cur)
}

func (c *Controller) deleteFlow(obj interface{}) {
	if flow, ok := obj.(*v1alpha1.Flow); ok {
		c.addFlow(flow)
		return
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
		return
	}
	flow, ok := tombstone.Obj.(*v1alpha1.Flow)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Pipe: %#v", tombstone.Obj))
		return
	}
	c.addFlow(flow)
}
