package flow

import (
	"fmt"
	"reflect"
	"time"

	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = v1alpha1.SchemeGroupVersion.WithKind("Flow")

const (
	// CodeSourceYAML defines code source config name
	CodeSourceYAML = "codesource.yaml"
	// ImageListYAML defines config name of image list
	ImageListYAML = "imagelist.yaml"
)

func (c *Controller) addFlow(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.flowQueue.Add(key)
}

func (c *Controller) updateFlow(old, cur interface{}) {
	oldFlow := old.(*v1alpha1.Flow)
	curFlow := cur.(*v1alpha1.Flow)

	if !reflect.DeepEqual(&oldFlow.Spec, &curFlow.Spec) {
		c.addFlow(cur)
	}
}

func (c *Controller) syncFlowHandler(key string) error {
	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing flow %q (%v)", key, time.Since(startTime))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	flow, err := c.flowLister.Flows(ns).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		delete(c.jobCaches, key)
		return nil
	}

	if err := c.syncFlow(flow); err != nil {
		return err
	}
	return nil
}

// syncFlow defines main process of CI/CD flow.
// It contains 4 steps:
//   1. TODO(liubog2008): update git cache, now just always use remote repo
//   2. create job to sync code from code source
//   3. create job of different stages
func (c *Controller) syncFlow(flow *v1alpha1.Flow) error {
	hash := DeepHashObject(flow.Spec)

	if err := c.syncConfigMap(flow); err != nil {
		return err
	}
	if err := c.syncPVC(flow); err != nil {
		return err
	}

	if err := c.GenerateJobs(flow, hash); err != nil {
		return err
	}

	oldJobs, curJobs, err := c.listAllJobsByFlow(flow, hash)
	if err != nil {
		return err
	}
	if err := c.cleanJobs(oldJobs); err != nil {
		return err
	}
	cur, next, err := c.curAndNextJob(curJobs, flow)
	if err != nil {
		return err
	}
	if cur == nil || isJobCompleted(cur) {
		nextStage := getJobStage(next)
		if nextStage == "" {
			if err := c.updateFlowStatus(flow, nextStage, "Running"); err != nil {
				return err
			}
			if err := c.tryRunNextJob(next); err != nil {
				return err
			}
			return nil
		}
		if err := c.updateFlowStatus(flow, nextStage, "Completed"); err != nil {
			return err
		}
		return nil
	}
	if isJobFailed(cur) {
		curStage := getJobStage(cur)
		if err := c.updateFlowStatus(flow, curStage, "Failed"); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (c *Controller) tryRunNextJob(job *batchv1.Job) error {
	if _, err := c.kubeclient.BatchV1().Jobs(job.Namespace).Create(job); err != nil {
		if errors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func (c *Controller) updateFlowStatus(flow *v1alpha1.Flow, stage, phase string) error {
	nf := flow.DeepCopy()
	nf.Status.Stage = stage
	nf.Status.Phase = phase
	if _, err := c.extclient.FlowV1alpha1().Flows(flow.Namespace).UpdateStatus(nf); err != nil {
		return err
	}
	return nil
}

func (c *Controller) curAndNextJob(jobs []batchv1.Job, flow *v1alpha1.Flow) (*batchv1.Job, *batchv1.Job, error) {
	codeReady, imageReady, appReady := splitStages(flow.Spec.Stages)
	stages := stageNameList(codeReady, imageReady, appReady)
	var curJob, nextJob *batchv1.Job
	cur := -1
	for i := len(stages) - 1; i >= 0; i-- {
		matched := false
		stage := stages[i]
		for _, job := range jobs {
			if getJobStage(&job) == stage {
				curJob = &job
				matched = true
				break
			}
		}
		if matched {
			cur = i
			break
		}
	}
	next := cur + 1
	if next != len(stages) {
		nextStage := stages[next]
		job, err := c.getJobWithStage(flow, nextStage)
		if err != nil {
			return nil, nil, err
		}
		nextJob = job
	}
	return curJob, nextJob, nil
}

func (c *Controller) listAllJobsByFlow(flow *v1alpha1.Flow, hash string) ([]batchv1.Job, []batchv1.Job, error) {
	opt := metav1.ListOptions{
		LabelSelector: flow.Spec.Selector.String(),
	}
	jobList, err := c.kubeclient.BatchV1().Jobs(flow.Namespace).List(opt)
	if err != nil {
		return nil, nil, err
	}
	oldJobs, curJobs := []batchv1.Job{}, []batchv1.Job{}
	for _, job := range jobList.Items {
		ref := metav1.GetControllerOf(&job)
		if ref.Kind == controllerKind.Kind && ref.UID == flow.UID {
			if getJobHash(&job) == hash {
				curJobs = append(curJobs, job)
			} else {
				oldJobs = append(oldJobs, job)
			}
		}
	}
	return oldJobs, curJobs, nil
}

func (c *Controller) cleanJobs(jobs []batchv1.Job) error {
	errList := []error{}
	for _, job := range jobs {
		if err := c.kubeclient.BatchV1().Jobs(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{
			Preconditions: &metav1.Preconditions{
				UID: &job.UID,
			},
		}); err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) != 0 {
		return fmt.Errorf("%v", errList)
	}
	return nil
}

func (c *Controller) getJobWithStage(flow *v1alpha1.Flow, stage string) (*batchv1.Job, error) {
	key, err := cache.MetaNamespaceKeyFunc(flow)
	if err != nil {
		return nil, err
	}
	jobCache, ok := c.jobCaches[key]
	if !ok {
		return nil, fmt.Errorf("can't hit job cache")
	}
	job, ok := jobCache.cache[stage]
	if !ok {
		return nil, fmt.Errorf("can't find job with stage %v", stage)
	}
	return job, nil
}

func (c *Controller) syncConfigMap(flow *v1alpha1.Flow) error {
	_, err := c.kubeclient.CoreV1().ConfigMaps(flow.Namespace).Get(flow.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		cm, err := c.generateConfigMap(flow)
		if err != nil {
			return err
		}
		if _, err := c.kubeclient.CoreV1().ConfigMaps(cm.Namespace).Create(cm); err != nil {
			return err
		}
	}
	return nil
}
func (c *Controller) syncPVC(flow *v1alpha1.Flow) error {
	_, err := c.kubeclient.CoreV1().PersistentVolumeClaims(flow.Namespace).Get(flow.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		pvc := c.generatePVC(flow)
		if _, err := c.kubeclient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(pvc); err != nil {
			return err
		}
	}
	return nil
}
