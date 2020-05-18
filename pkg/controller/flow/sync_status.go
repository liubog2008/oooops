package flow

import (
	"fmt"
	"reflect"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

const (
	// gitSuffix is name suffix of job which will init volume by git
	gitSuffix = "git"
)

func (c *Controller) syncFlowStatus(flow *v1alpha1.Flow, jobMap map[string]*batchv1.Job, pvc *corev1.PersistentVolumeClaim) (*v1alpha1.Flow, error) {
	status, err := c.generateFlowStatus(flow, jobMap, pvc)
	if err != nil {
		return nil, err
	}

	// TODO(liubog2008): optimize this function
	if reflect.DeepEqual(status, &flow.Status) {
		return nil, nil
	}

	updating := v1alpha1.Flow{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       flow.Namespace,
			Name:            flow.Name,
			ResourceVersion: flow.ResourceVersion,
		},
		Status: *status,
	}

	updated, err := c.extClient.MarioV1alpha1().Flows(updating.Namespace).UpdateStatus(&updating)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (c *Controller) calculateStageStatus(stages []v1alpha1.Stage, jobMap map[string]*batchv1.Job) ([]v1alpha1.StageStatus, error) {
	stageStatuses := []v1alpha1.StageStatus{}
	missingCount := 0
	for i := range stages {
		stage := &stages[i]

		job, ok := jobMap[stage.Name]
		if !ok {
			missingCount++
			continue
		}

		if missingCount != 0 {
			for k := 0; k < missingCount; k++ {
				// job is not found
				stageStatuses = append(stageStatuses, v1alpha1.StageStatus{
					Phase: v1alpha1.StageJobMissing,
				})
			}
			missingCount = 0
		}

		jobPhase := v1alpha1.StageJobRunning

		if IsJobComplete(job) {
			jobPhase = v1alpha1.StageJobComplete
		}

		if IsJobFailed(job) {
			jobPhase = v1alpha1.StageJobFailed
		}

		stageStatuses = append(stageStatuses, v1alpha1.StageStatus{
			Job:   job.Name,
			Phase: jobPhase,
		})
	}
	return stageStatuses, nil
}

func (c *Controller) generateFlowStatus(flow *v1alpha1.Flow, jobMap map[string]*batchv1.Job, pvc *corev1.PersistentVolumeClaim) (*v1alpha1.FlowStatus, error) {
	status := v1alpha1.FlowStatus{}

	pvcCond := generateGitVolumeCondition(pvc)
	status.Conditions = append(status.Conditions, *pvcCond)

	gitJob := jobMap["git"]
	marioJob := jobMap["mario"]

	marioCond := generateMarioCondition(flow, gitJob, marioJob)
	status.Conditions = append(status.Conditions, *marioCond)

	stageStatuses, err := c.calculateStageStatus(flow.Spec.Stages, jobMap)
	if err != nil {
		return nil, err
	}
	status.StageStatuses = stageStatuses
	length := len(stageStatuses)

	if length == 0 {
		status.Phase = v1alpha1.FlowPending
		return &status, nil
	}

	lastPhase := stageStatuses[length-1].Phase

	if lastPhase == v1alpha1.StageJobFailed {
		status.Phase = v1alpha1.FlowFailed
		return &status, nil
	}

	if lastPhase == v1alpha1.StageJobComplete &&
		len(status.StageStatuses) == len(flow.Spec.Stages) {
		status.Phase = v1alpha1.FlowSucceed
		return &status, nil
	}

	status.Phase = v1alpha1.FlowRunning
	return &status, nil
}

func generateGitVolumeCondition(pvc *corev1.PersistentVolumeClaim) *v1alpha1.FlowCondition {
	if pvc == nil {
		return NewFlowCondition(
			v1alpha1.FlowGitVolumeReady,
			corev1.ConditionFalse,
			v1alpha1.FlowReasonGitVolumeClaiming,
			"PersistentVolumeClaim is not created",
		)
	}
	switch pvc.Status.Phase {
	case corev1.ClaimPending:
		return NewFlowCondition(
			v1alpha1.FlowGitVolumeReady,
			corev1.ConditionFalse,
			v1alpha1.FlowReasonGitVolumePending,
			"PersistentVolumeClaim is in pending phase",
		)
	case corev1.ClaimBound:
		return NewFlowCondition(
			v1alpha1.FlowGitVolumeReady,
			corev1.ConditionTrue,
			v1alpha1.FlowReasonGitVolumeBound,
			"PersistentVolumeClaim is bound",
		)
	case corev1.ClaimLost:
		return NewFlowCondition(
			v1alpha1.FlowGitVolumeReady,
			corev1.ConditionFalse,
			v1alpha1.FlowReasonGitVolumeLost,
			"Volume of the PersistentVolumeClaim is lost",
		)
	}
	return NewFlowCondition(
		v1alpha1.FlowGitVolumeReady,
		corev1.ConditionUnknown,
		v1alpha1.FlowReasonGitVolumeUnknown,
		"PersistentVolumeClaim status is unknown",
	)
}

func generateMarioCondition(flow *v1alpha1.Flow, gitJob, marioJob *batchv1.Job) *v1alpha1.FlowCondition {
	if flow.Spec.Mario != nil {
		return NewFlowCondition(
			v1alpha1.FlowMarioReady,
			corev1.ConditionTrue,
			v1alpha1.FlowReasonMarioReady,
			"Mario is ready",
		)
	}

	if marioJob != nil {
		if IsJobFailed(marioJob) {
			return NewFlowCondition(
				v1alpha1.FlowMarioReady,
				corev1.ConditionFalse,
				v1alpha1.FlowReasonMarioFailed,
				fmt.Sprintf("Mario job %v is failed", marioJob.Name),
			)
		}

		return NewFlowCondition(
			v1alpha1.FlowMarioReady,
			corev1.ConditionFalse,
			v1alpha1.FlowReasonMarioPending,
			"Waiting for mario job to init mario",
		)
	}

	if gitJob != nil {
		if IsJobFailed(gitJob) {
			return NewFlowCondition(
				v1alpha1.FlowMarioReady,
				corev1.ConditionFalse,
				v1alpha1.FlowReasonGitFailed,
				fmt.Sprintf("Git job %v is failed", gitJob.Name),
			)
		}

		if IsJobComplete(gitJob) {
			return NewFlowCondition(
				v1alpha1.FlowMarioReady,
				corev1.ConditionFalse,
				v1alpha1.FlowReasonMarioPending,
				"Waiting for mario job to init mario",
			)
		}

		return NewFlowCondition(
			v1alpha1.FlowMarioReady,
			corev1.ConditionFalse,
			v1alpha1.FlowReasonGitPending,
			"Waiting for git job to fetch code",
		)

	}

	return NewFlowCondition(
		v1alpha1.FlowMarioReady,
		corev1.ConditionFalse,
		v1alpha1.FlowReasonGitPending,
		"Waiting for git job to fetch code",
	)
}
