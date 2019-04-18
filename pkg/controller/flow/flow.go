package flow

import (
	"reflect"
	"time"

	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
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
	jobs, err := c.kubeclient.BatchV1().Jobs(flow.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) listAllJobsByFlow(flow *v1alpha1.Flow) ([]batchv1.Job, error) {
	opt := metav1.ListOptions{
		LabelSelector: flow.Spec.Selector,
	}
	jobs, err := c.kubeclient.BatchV1().Jobs(flow.Namespace).List(opt)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs {
	}
}

func (c *Controller) GenerateJobs(flow *v1alpha1.Flow) ([]batchv1.Job, error) {
}

func (c *Controller) generateSCMJob(source *v1alpha1.CodeSource) (*batchv1.Job, error) {
	// TODO(liubog2008): support more code source
	return c.generateGitJob(&source.GitSource)
}
