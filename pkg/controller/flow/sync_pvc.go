package flow

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

func (c *Controller) syncPVC(flow *v1alpha1.Flow, pvc *corev1.PersistentVolumeClaim) error {
	if pvc != nil {
		if !metav1.IsControlledBy(pvc, flow) {
			return fmt.Errorf("can't create pvc %s/%s, it exists and is not controlled by flow")
		}

		// NOTE(liubog2008): handle updation of pvc?
		return nil

	}
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)
	pvc = flow.Spec.Git.VolumeClaimTemplate.DeepCopy()

	pvc.Name = flow.Name
	pvc.Namespace = flow.Namespace
	pvc.OwnerReferences = append(pvc.OwnerReferences, *owner)

	if pvc.Labels == nil {
		pvc.Labels = flow.Spec.Selector.MatchLabels
	} else {
		for k, v := range flow.Spec.Selector.MatchLabels {
			pvc.Labels[k] = v
		}
	}

	if _, err := c.kubeClient.CoreV1().PersistentVolumeClaims(flow.Namespace).Create(pvc); err != nil {
		return err
	}

	return nil

}
