package flow

import (
	"fmt"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

const (
	gitScriptName       = "git.sh"
	gitScriptPath       = "/app"
	gitScriptVolumeName = "git-script"
)

func (c *Controller) syncJob(flow *v1alpha1.Flow, jobMap map[string]*batchv1.Job) error {
	// mario is not set, generate init job
	if flow.Spec.Mario == nil {
		gitJob, gitOK := jobMap[v1alpha1.FlowStageGit]
		_, marioOK := jobMap[v1alpha1.FlowStageMario]

		if !gitOK && !marioOK {
			gitJob = c.generateGitJob(flow)
			if _, err := c.kubeClient.BatchV1().Jobs(flow.Namespace).Create(gitJob); err != nil {
				return err
			}
			return nil
		}
		if !marioOK && IsJobComplete(gitJob) {
			marioJob := c.generateMarioJob(flow)
			if _, err := c.kubeClient.BatchV1().Jobs(flow.Namespace).Create(marioJob); err != nil {
				return err
			}
		}
		return nil
	}

	job, err := c.generateNextJob(flow, jobMap)
	if err != nil {
		return err
	}

	if job == nil {
		return nil
	}

	if _, err := c.kubeClient.BatchV1().Jobs(flow.Namespace).Create(job); err != nil {
		return err
	}

	return nil
}

func (c *Controller) generateNextJob(flow *v1alpha1.Flow, jobMap map[string]*batchv1.Job) (*batchv1.Job, error) {
	last := -1
	for i := range flow.Spec.Stages {
		stage := &flow.Spec.Stages[i]

		_, ok := jobMap[v1alpha1.UserJobPrefix+stage.Name]
		if ok {
			last = i
		}
	}
	// all jobs have been generated
	if last == len(flow.Spec.Stages)-1 {
		return nil, nil
	}

	if last != -1 {
		lastStage := &flow.Spec.Stages[last]

		// existence has checked
		lastJob := jobMap[v1alpha1.UserJobPrefix+lastStage.Name]

		if !IsJobComplete(lastJob) {
			return nil, nil
		}
	}

	curIndex := last + 1

	// NOTE: len of status stages are always less than len of spec stages
	// TODO(liubog2008): add test case to test it
	// only when last stage has been completed
	// next job will be generated
	job, err := c.generateActionJob(flow, curIndex)

	return job, err
}

func (c *Controller) generateActionJob(flow *v1alpha1.Flow, stageIndex int) (*batchv1.Job, error) {
	stage := flow.Spec.Stages[stageIndex]
	mario := flow.Spec.Mario
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)

	for i := range mario.Spec.Actions {
		action := &mario.Spec.Actions[i]
		if action.Name != stage.Action {
			continue
		}

		labels := map[string]string{}
		for k, v := range flow.Spec.Selector.MatchLabels {
			labels[k] = v
		}
		for k, v := range mario.Labels {
			_, ok := labels[k]
			if !ok {
				labels[k] = v
			}
		}

		version := flow.Spec.Git.Ref

		cs, err := constructContainers(action, version)
		if err != nil {
			return nil, err
		}

		job := batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nameJoin(flow.Name, "user", stage.Name),
				Namespace: flow.Namespace,
				Labels:    labels,
				OwnerReferences: []metav1.OwnerReference{
					*owner,
				},
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						Containers:    cs,
						Volumes: []corev1.Volume{
							{
								Name: gitRootVolumeName,
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: flow.Name,
									},
								},
							},
						},
					},
				},
			},
		}

		return &job, nil
	}

	return nil, nil
}

func constructContainers(action *v1alpha1.MarioAction, version string) ([]corev1.Container, error) {
	// TODO(liubog2008): support action importing
	if action.Template == nil {
		return nil, fmt.Errorf("no action template")
	}
	env := make([]corev1.EnvVar, 0, len(action.Env)+1)
	env = append(env, corev1.EnvVar{
		Name:  action.Template.Version.EnvName,
		Value: version,
	})

	for i := range action.Env {
		e := &action.Env[i]
		if e.Name == action.Template.Version.EnvName {
			return nil, fmt.Errorf("set an env whose name is confict with version env")
		}
		env = append(env, corev1.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}

	c := corev1.Container{
		Name:       action.Name,
		Image:      action.Template.Image,
		Command:    action.Template.Command,
		Args:       action.Template.Args,
		WorkingDir: action.Template.WorkingDir,

		Env: env,

		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      gitRootVolumeName,
				MountPath: action.Template.WorkingDir,
			},
		},
	}

	return []corev1.Container{c}, nil
}

func (c *Controller) generateGitJob(flow *v1alpha1.Flow) *batchv1.Job {
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)

	fileMode := int32(0755)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name + "-git",
			Namespace: flow.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*owner,
			},
			Labels: flow.Spec.Selector.MatchLabels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "git",
							Image: c.gitImage,
							Command: []string{
								filepath.Join(gitScriptPath, gitScriptName),
							},
							WorkingDir: marioWorkingDir,

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      gitRootVolumeName,
									MountPath: marioWorkingDir,
								},
								{
									Name:      gitScriptVolumeName,
									MountPath: gitScriptPath,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: gitRootVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: flow.Name,
								},
							},
						},
						{
							Name: gitScriptVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{

										Name: flow.Name,
									},
									DefaultMode: &fileMode,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *Controller) generateMarioJob(flow *v1alpha1.Flow) *batchv1.Job {
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)

	command := []string{
		"/app/mario",
		"--remote",
		flow.Spec.Git.Repo,
		"--ref",
		flow.Spec.Git.Ref,
		"--addr",
		":8080",
		"--token",
		"test",
	}

	labels := map[string]string{}
	for k, v := range flow.Spec.Selector.MatchLabels {
		labels[k] = v
	}

	labels[v1alpha1.DefaultFlowStageLabelKey] = v1alpha1.FlowStageMario

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name + "-mario",
			Namespace: flow.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*owner,
			},
			Labels: labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "mario",
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:       "mario",
							Image:      c.marioImage,
							Command:    command,
							WorkingDir: marioWorkingDir,

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      gitRootVolumeName,
									MountPath: marioWorkingDir,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.FromInt(8080),
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: gitRootVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: flow.Name,
								},
							},
						},
					},
				},
			},
		},
	}
}
