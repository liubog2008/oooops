package flow

import (
	"fmt"
	"reflect"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

func (c *Controller) addJob(obj interface{}) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("obj is not Job: %v", obj))
		return
	}

	flow := c.getFlowFromJob(job)
	if flow == nil {
		return
	}

	klog.Infof("enqueue flow %s/%s by job %s", flow.Namespace, flow.Name, job.Name)

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

	refChanged := !reflect.DeepEqual(oldRef, curRef)

	if oldRef != nil && refChanged {
		if oldFlow := c.getFlowFromRef(oldJob.Namespace, oldRef); oldFlow != nil {
			klog.Infof("enqueue flow %s/%s by old job %s", oldFlow.Namespace, oldFlow.Name, oldJob.Name)
			c.addFlow(oldFlow)
		}
	}
	if curRef != nil {
		if curFlow := c.getFlowFromRef(curJob.Namespace, curRef); curFlow != nil {
			klog.Infof("enqueue flow %s/%s by current job %s", curFlow.Namespace, curFlow.Name, curJob.Name)
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
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Job: %#v", tombstone.Obj))
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
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Flow: %#v", tombstone.Obj))
		return
	}
	c.addFlow(flow)
}

func (c *Controller) addPod(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("obj is not Pod: %v", obj))
		return
	}
	// ignore unready pod
	if !IsPodReady(pod) {
		return
	}

	flow := c.getFlowFromPod(pod)
	if flow == nil {
		return
	}

	klog.Infof("enqueue flow %s/%s by pod %s", flow.Namespace, flow.Name, pod.Name)

	c.addFlow(flow)
}

func (c *Controller) updatePod(old, cur interface{}) {
	oldPod, ok1 := old.(*corev1.Pod)
	curPod, ok2 := cur.(*corev1.Pod)
	if !ok1 || !ok2 {
		utilruntime.HandleError(fmt.Errorf("either old or cur is not Pod: %v, %v", old, cur))
		return
	}
	if oldPod.ResourceVersion == curPod.ResourceVersion {
		return
	}

	c.addPod(curPod)
}

func (c *Controller) deletePod(obj interface{}) {
	if pod, ok := obj.(*corev1.Pod); ok {
		c.addPod(pod)
		return
	}
	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
		return
	}
	pod, ok := tombstone.Obj.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Pod: %#v", tombstone.Obj))
		return
	}
	c.addPod(pod)
}
