package flow

import (
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

func (c *Controller) getFlowFromRef(ns string, ref *metav1.OwnerReference) *v1alpha1.Flow {
	if ref.Kind != c.Kind {
		return nil
	}

	flow, err := c.flowLister.Flows(ns).Get(ref.Name)
	if err != nil {
		return nil
	}

	if flow.UID != ref.UID {
		return nil
	}
	return flow
}

// NewFlowCondition returns a condition of flow
func NewFlowCondition(t v1alpha1.FlowConditionType, status corev1.ConditionStatus, reason, message string) *v1alpha1.FlowCondition {
	return &v1alpha1.FlowCondition{
		Type:               t,
		Status:             status,
		LastProbeTime:      metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func IsJobComplete(job *batchv1.Job) bool {
	s, ok := getJobCondition(job, batchv1.JobComplete)
	if !ok {
		return false
	}
	return s == corev1.ConditionTrue
}

func IsJobFailed(job *batchv1.Job) bool {
	s, ok := getJobCondition(job, batchv1.JobFailed)
	if !ok {
		return false
	}
	return s == corev1.ConditionTrue
}

func getJobCondition(job *batchv1.Job, t batchv1.JobConditionType) (corev1.ConditionStatus, bool) {
	for i := range job.Status.Conditions {
		c := &job.Status.Conditions[i]

		if c.Type == t {
			return c.Status, true
		}
	}
	return "", false
}

func (c *Controller) cleanInvalidJob(message, namespace, name string) {
	klog.Warningf(message)
	if err := c.kubeClient.BatchV1().Jobs(namespace).Delete(name, nil); err != nil {
		klog.Warningf("clean job %v/%v failed: %v", namespace, name, err)
	}
}

func nameJoin(parts ...string) string {
	return strings.Join(parts, "-")
}
