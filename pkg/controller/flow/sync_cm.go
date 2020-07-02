package flow

import (
	"bytes"
	"fmt"
	"html/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
)

const (
	gitScriptTemplateContent = `#!/bin/sh
git init
git remote add origin {{ .repo }}

git fetch origin {{ .ref }} --depth=1

git reset --hard FETCH_HEAD`
)

var (
	gitScriptTemplate = template.Must(
		template.New("git-script").Parse(gitScriptTemplateContent),
	)
)

func (c *Controller) syncConfigMap(flow *v1alpha1.Flow, cm *corev1.ConfigMap) error {
	if cm != nil {
		if !metav1.IsControlledBy(cm, flow) {
			// TODO(liubog2008): fix configmap name conflict
			return fmt.Errorf("can't create cm %s/%s, it exists and is not controlled by flow", cm.Namespace, cm.Name)
		}

		// NOTE(liubog2008): handle updation of cm?
		return nil
	}

	gitScript, err := c.generateGitScript(flow)
	if err != nil {
		return err
	}

	if _, err := c.kubeClient.CoreV1().ConfigMaps(flow.Namespace).Create(gitScript); err != nil {
		return err
	}

	return nil
}

func (c *Controller) generateGitScript(flow *v1alpha1.Flow) (*corev1.ConfigMap, error) {
	owner := metav1.NewControllerRef(flow, c.GroupVersionKind)
	buf := bytes.Buffer{}

	if err := gitScriptTemplate.Execute(&buf, map[string]string{
		"repo": flow.Spec.Git.Repo,
		"ref":  flow.Spec.Git.Ref,
	}); err != nil {
		return nil, err
	}

	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flow.Name,
			Namespace: flow.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*owner,
			},
		},
		Data: map[string]string{
			gitScriptName: buf.String(),
		},
	}

	return &cm, nil
}
