package flow

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

func (c *Controller) syncJob(flow *v1alpha1.Flow, jobMap map[string]*batchv1.Job) error {
	// mario is not set, generate init job
	if flow.Spec.Mario == nil {
		gitJob, gitOK := jobMap["git"]
		marioJob, marioOK := jobMap["mario"]

		if !gitOK && !marioOK {
			gitJob = c.generateGitJob(flow)
			if _, err := c.kubeClient.BatchV1().Jobs(flow.Namespace).Create(gitJob); err != nil {
				return err
			}
		}
		if !marioOK && IsJobComplete(gitJob) {
			marioJob = c.generateMarioJob(flow)
			if _, err := c.kubeClient.BatchV1().Jobs(flow.Namespace).Create(marioJob); err != nil {
				return err
			}
		}
		return nil
	}

	job := c.generateNextJob(flow)
	if job == nil {
		return nil
	}

	if _, err := c.kubeClient.BatchV1().Jobs(flow.Namespace).Create(job); err != nil {
		return err
	}

	return nil
}

func (c *Controller) generateNextJob(flow *v1alpha1.Flow) *batchv1.Job {
	// all jobs have been generated
	if len(flow.Status.StageStatuses) == len(flow.Spec.Stages) {
		return nil
	}

	index := len(flow.Status.StageStatuses)
	if index != 0 {
		// last stage is not completed
		if flow.Status.StageStatuses[index-1].Phase != v1alpha1.StageJobComplete {
			return nil
		}
	}

	// NOTE: len of status stages are always less than len of spec stages
	// TODO(liubog2008): add test case to test it
	// only when last stage has been completed
	// next job will be generated
	job := c.generateActionJob(flow, index)

	return job
}

func (c *Controller) generateActionJob(flow *v1alpha1.Flow, stageIndex int) *batchv1.Job {
	stage := flow.Spec.Stages[stageIndex]
	mario := flow.Spec.Mario
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)

	for i := range mario.Spec.Actions {
		action := &mario.Spec.Actions[i]
		if action.Name != stage.Action {
			continue
		}

		labels := map[string]string{}
		for k, v := range mario.Labels {
			labels[k] = v
		}

		version := flow.Spec.Git.Ref

		cs := constructContainers(action, version)

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

		return &job
	}

	return nil
}

func constructContainers(action *v1alpha1.Action, version string) []corev1.Container {
	cs := action.Template.Spec.Containers
	kcs := make([]corev1.Container, 0, len(cs))
	for i := range cs {
		c := &cs[i]
		kcs = append(kcs, corev1.Container{
			Name:       c.Name,
			Image:      c.Image,
			Command:    c.Command,
			Args:       c.Args,
			WorkingDir: action.WorkingDir,

			Env: []corev1.EnvVar{
				{
					Name:  action.Version.EnvName,
					Value: version,
				},
			},

			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      gitRootVolumeName,
					MountPath: action.WorkingDir,
				},
			},
		})
	}

	return kcs
}

func (c *Controller) generateGitJob(flow *v1alpha1.Flow) *batchv1.Job {
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)

	// FIXME(liubog2008): make it dynamic and can be changed by config
	command := []string{
		"git",
		"clone",
		flow.Spec.Git.Repo + "@" + flow.Spec.Git.Ref,
		"--depth",
		"1",
	}

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
							Name:       "git",
							Image:      c.gitImage,
							Command:    command,
							WorkingDir: marioWorkingDir,

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      gitRootVolumeName,
									MountPath: marioWorkingDir,
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

func (c *Controller) generateMarioJob(flow *v1alpha1.Flow) *batchv1.Job {
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name + "-mario",
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
							Name:       "mario",
							Image:      c.marioImage,
							WorkingDir: marioWorkingDir,

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      gitRootVolumeName,
									MountPath: marioWorkingDir,
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
