package flow

import (
	"fmt"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

func (c *Controller) syncFlow(key string) error {
	startTime := time.Now()

	defer func() {
		klog.V(4).Infof("Finished syncing flow %q. (%v)", key, time.Since(startTime))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	flow, err := c.flowLister.Flows(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	selector, err := metav1.LabelSelectorAsSelector(flow.Spec.Selector)
	if err != nil {
		return fmt.Errorf("converting flow selector error: %v", err)
	}

	jobs, err := c.jobLister.Jobs(ns).List(selector)
	if err != nil {
		return err
	}

	klog.Infof("number of jobs selected by selector %v: %v", selector, len(jobs))

	pvc, err := c.pvcLister.PersistentVolumeClaims(ns).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	jobMap := c.calculateJobMap(flow, jobs)

	updated, err := c.syncFlowStatus(flow, jobMap, pvc)
	if err != nil {
		return err
	}

	// status is not updated
	if updated == nil {
		return nil
	}

	if err := c.syncPVC(flow, pvc); err != nil {
		return err
	}

	if err := c.syncJob(updated, jobMap); err != nil {
		return err
	}

	return nil
}

// calculateJobMap calculate job map of flow and returns the map and
// current stage index
func (c *Controller) calculateJobMap(flow *v1alpha1.Flow, jobs []*batchv1.Job) map[string]*batchv1.Job {
	jobMap := map[string]*batchv1.Job{}
	for _, job := range jobs {
		klog.Infof("check job %v", job.Name)
		if !metav1.IsControlledBy(job, flow) {
			// ignore jobs which is not controlled by the flow
			continue
		}
		stage := strings.TrimPrefix(job.Name, flow.Name+"-")

		switch stage {
		case "git":
		case "mario":
		default:
			if !strings.HasPrefix(stage, v1alpha1.UserJobPrefix) {
				msg := fmt.Sprintf("unsupported system stage %s of job(%s/%s)", stage, job.Namespace, job.Name)
				c.cleanInvalidJob(msg, job.Namespace, job.Name)
				continue
			}
		}
		other, ok := jobMap[stage]
		if ok {
			older := job
			if !older.CreationTimestamp.Before(&other.CreationTimestamp) {
				older = other
				jobMap[stage] = job
			}

			msg := fmt.Sprintf(
				"job (%s/%s) has same stage %s with another job(%s/%s), clean the older one(%s/%s)",
				job.Namespace, job.Name,
				stage,
				other.Namespace, other.Name,
				older.Namespace, older.Name,
			)
			c.cleanInvalidJob(msg, older.Namespace, older.Name)
			continue
		}
		jobMap[stage] = job
	}

	return jobMap
}
