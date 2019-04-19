package flow

import (
	"fmt"
	"hash/fnv"

	"github.com/davecgh/go-spew/spew"
	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
)

// splitStages splits stages by its when field
func splitStages(stages []v1alpha1.Stage) (codeReady []v1alpha1.Stage, imageReady []v1alpha1.Stage, appReady []v1alpha1.Stage) {
	for _, stage := range stages {
		switch stage.When {
		case v1alpha1.WhenCodeReady:
			codeReady = append(codeReady, stage)
		case v1alpha1.WhenImageReady:
			imageReady = append(imageReady, stage)
		case v1alpha1.WhenApplicationReady:
			appReady = append(appReady, stage)
		}
	}
	return codeReady, imageReady, appReady
}

func stageNameList(codeReady, imageReady, appReady []v1alpha1.Stage) []string {
	names := []string{v1alpha1.StageSCM}
	for _, s := range codeReady {
		names = append(names, s.Name)
	}
	names = append(names, v1alpha1.StageImage)
	for _, s := range imageReady {
		names = append(names, s.Name)
	}
	names = append(names, v1alpha1.StageDeploy)
	for _, s := range appReady {
		names = append(names, s.Name)
	}
	return names
}

func getJobStage(job *batchv1.Job) string {
	if job.Labels == nil {
		return ""
	}
	s, ok := job.Labels[v1alpha1.LabelStage]
	if !ok {
		return ""
	}
	return s
}

func getJobHash(job *batchv1.Job) string {
	if job.Labels == nil {
		return ""
	}
	s, ok := job.Labels[v1alpha1.LabelFlowHash]
	if !ok {
		return ""
	}
	return s
}

// DeepHashObject hashes an object
func DeepHashObject(objectToWrite interface{}) string {
	hasher := fnv.New64a()
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", objectToWrite)
	return fmt.Sprint(hasher.Sum64())
}
