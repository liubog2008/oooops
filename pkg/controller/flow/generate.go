package flow

import (
	"fmt"

	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	yaml "gopkg.in/yaml.v2"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) GenerateJobs(flow *v1alpha1.Flow, hash string) error {
	key, err := cache.MetaNamespaceKeyFunc(flow)
	if err != nil {
		return err
	}
	jobCache, ok := c.jobCaches[key]
	if ok && jobCache.hash == hash {
		return nil
	}
	jobMap, err := c.generateJobs(flow, hash)
	if err != nil {
		return err
	}
	c.jobCaches[key] = JobCache{
		hash:  hash,
		cache: jobMap,
	}
	return nil
}

func (c *Controller) generateJobs(flow *v1alpha1.Flow, hash string) (map[string]*batchv1.Job, error) {
	jobMap := map[string]*batchv1.Job{}
	jobMap[v1alpha1.StageSCM] = c.generateSCMJob(flow, hash)
	jobMap[v1alpha1.StageImage] = c.generateImageJob(flow, hash)
	// jobMap[v1alpha1.StageDeploy] = c.generateDeploy(flow)

	return jobMap, nil
}

func (c *Controller) generateSCMJob(flow *v1alpha1.Flow, hash string) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name + "-scm",
			Namespace: flow.Namespace,
			Labels: map[string]string{
				v1alpha1.LabelStage:    v1alpha1.StageSCM,
				v1alpha1.LabelFlowHash: hash,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(flow, controllerKind),
			},
		},
	}
	// FIXME(liubog2008): fix label conflict if selector contains
	// v1alpha1.LabelStage
	for labelKey, labelValue := range flow.Spec.Selector.MatchLabels {
		job.Labels[labelKey] = labelValue
	}

	job.Spec = batchv1.JobSpec{
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: job.Labels,
			},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				Containers: []v1.Container{
					{
						Name:  "scm",
						Image: c.scmImage,
						Args: []string{
							"git",
							"-d",
							"/gitrepo",
							"/config/" + CodeSourceYAML,
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "gitrepo",
								MountPath: "/gitrepo",
							},
							{
								Name:      "config",
								MountPath: "/config/" + CodeSourceYAML,
								SubPath:   CodeSourceYAML,
							},
						},
					},
				},
				ImagePullSecrets: []v1.LocalObjectReference{
					{
						Name: "aliyun",
					},
				},
				Volumes: []v1.Volume{
					{
						Name: "gitrepo",
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: flow.Name,
							},
						},
					},
					{
						Name: "config",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: flow.Name,
								},
								Items: []v1.KeyToPath{
									{
										Key:  CodeSourceYAML,
										Path: CodeSourceYAML,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return job
}

func (c *Controller) generateImageJob(flow *v1alpha1.Flow, hash string) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name + "-image",
			Namespace: flow.Namespace,
			Labels: map[string]string{
				v1alpha1.LabelStage:    v1alpha1.StageImage,
				v1alpha1.LabelFlowHash: hash,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(flow, controllerKind),
			},
		},
	}
	// FIXME(liubog2008): fix label conflict if selector contains
	// v1alpha1.LabelStage
	for labelKey, labelValue := range flow.Spec.Selector.MatchLabels {
		job.Labels[labelKey] = labelValue
	}

	job.Spec = batchv1.JobSpec{
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: job.Labels,
			},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				Containers: []v1.Container{
					{
						Name:  "image",
						Image: c.scmImage,
						Args: []string{
							"-d",
							"/gitrepo",
							"/config/" + ImageListYAML,
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "gitrepo",
								MountPath: "/gitrepo",
							},
							{
								Name:      "config",
								MountPath: "/config/" + ImageListYAML,
								SubPath:   ImageListYAML,
							},
						},
					},
				},
				ImagePullSecrets: []v1.LocalObjectReference{
					{
						Name: "aliyun",
					},
				},
				Volumes: []v1.Volume{
					{
						Name: "gitrepo",
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: flow.Name,
							},
						},
					},
					{
						Name: "config",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: flow.Name,
								},
								Items: []v1.KeyToPath{
									{
										Key:  ImageListYAML,
										Path: ImageListYAML,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return job
}

func (c *Controller) generateConfigMap(flow *v1alpha1.Flow) (*v1.ConfigMap, error) {
	codeSourceYAMLFile, err := c.generateCodeSourceYAML(flow)
	if err != nil {
		return nil, fmt.Errorf("can't generate code source config yaml")
	}
	imagesYAMLFile, err := c.generateImagesYAML(flow)
	if err != nil {
		return nil, fmt.Errorf("can't generate images config yaml")
	}
	cm := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name,
			Namespace: flow.Namespace,
		},
		Data: map[string]string{
			CodeSourceYAML: codeSourceYAMLFile,
			ImageListYAML:  imagesYAMLFile,
		},
	}
	for labelKey, labelValue := range flow.Spec.Selector.MatchLabels {
		cm.Labels[labelKey] = labelValue
	}

	cm.OwnerReferences = append(cm.OwnerReferences, *metav1.NewControllerRef(flow, controllerKind))
	return &cm, nil
}

func (c *Controller) generatePVC(flow *v1alpha1.Flow) *v1.PersistentVolumeClaim {
	pvc := flow.Spec.VolumeClaimTemplate.DeepCopy()
	pvc.Name = flow.Name
	pvc.Namespace = flow.Namespace
	pvc.OwnerReferences = append(pvc.OwnerReferences, *metav1.NewControllerRef(flow, controllerKind))
	for labelKey, labelValue := range flow.Spec.Selector.MatchLabels {
		pvc.Labels[labelKey] = labelValue
	}

	return pvc
}

func (c *Controller) generateCodeSourceYAML(flow *v1alpha1.Flow) (string, error) {
	body, err := yaml.Marshal(&flow.Spec.Source)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Controller) generateImagesYAML(flow *v1alpha1.Flow) (string, error) {
	body, err := yaml.Marshal(&flow.Spec.Images)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
