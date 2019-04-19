package flow

import (
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func (c *Controller) addJob(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.jobQueue.Add(key)
}

func (c *Controller) updateJob(old, cur interface{}) {
	curJob := cur.(*batchv1.Job)

	c.jobQueue.Add(curJob)
}

func (c *Controller) deleteJob(obj interface{}) {
	job, ok := obj.(*batchv1.Job)

	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %+v", obj))
			return
		}
		job, ok = tombstone.Obj.(*batchv1.Job)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a job %#v", obj))
			return
		}
	}
	c.jobQueue.Add(job)
}

func (c *Controller) syncJobHandler(key string) error {
	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing flow job %q (%v)", key, time.Since(startTime))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	job, err := c.jobLister.Jobs(ns).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	if err := c.syncJob(job); err != nil {
		return err
	}
	return nil
}

func (c *Controller) syncJob(job *batchv1.Job) error {
	stage := getJobStage(job)
	if stage == "" {
		return nil
	}
	ref := metav1.GetControllerOf(job)
	if ref == nil {
		return nil
	}
	if ref.Kind != controllerKind.Kind {
		return nil
	}
	// TODO(liubog2008): check condition before get flow
	flow, err := c.extclient.FlowV1alpha1().Flows(job.Namespace).Get(ref.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if stage != flow.Status.Stage {
		return nil
	}
	if isJobFinished(job) {
		c.addFlow(flow)
	}
	return nil
}

func checkJobCondition(job *batchv1.Job, ct batchv1.JobConditionType, status v1.ConditionStatus) bool {
	for _, c := range job.Status.Conditions {
		if c.Type == ct && c.Status == status {
			return true
		}
	}
	return false
}

func isJobFinished(job *batchv1.Job) bool {
	for _, c := range job.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobCompleted(job *batchv1.Job) bool {
	for _, c := range job.Status.Conditions {
		if c.Type == batchv1.JobComplete && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobFailed(job *batchv1.Job) bool {
	for _, c := range job.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
