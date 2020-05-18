// Package config defines config of operator
package config

import (
	"k8s.io/client-go/informers"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/liubog2008/oooops/pkg/client/clientset"
	extinformers "github.com/liubog2008/oooops/pkg/client/informers"
	marioinformers "github.com/liubog2008/oooops/pkg/client/informers/mario/v1alpha1"
)

// Config defines config of operator
type Config struct {
	KubeClient kubernetes.Interface

	ExtClient clientset.Interface

	// KubeInformerFactory defines informer factory
	KubeInformerFactory informers.SharedInformerFactory

	// ExtInformerFactory defines extension informer factory
	ExtInformerFactory extinformers.SharedInformerFactory

	EventInformer marioinformers.EventInformer

	PipeInformer marioinformers.PipeInformer

	FlowInformer marioinformers.FlowInformer

	JobInformer batchinformers.JobInformer

	PVCInformer coreinformers.PersistentVolumeClaimInformer
}
