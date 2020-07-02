// Package options defines options  of operator
package options

import (
	"fmt"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/liubog2008/oooops/cmd/operator/app/config"
	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/client/clientset"
	extinformers "github.com/liubog2008/oooops/pkg/client/informers"
)

// Options defines running options of operator
type Options struct {
	Kubeconfig string

	Namespace string
}

// NewOptions returns new running options
func NewOptions() (*Options, error) {
	opt := &Options{
		Kubeconfig: "",
		Namespace:  "default",
	}

	return opt, nil
}

// AddFlags adds flags for operator options
func (opt *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&opt.Kubeconfig, "kubeconfig", opt.Kubeconfig,
		"kubeconfig for cluster")
	fs.StringVar(&opt.Namespace, "namespace", opt.Namespace,
		"namespace which operator watches, if empty, all namespaces will be watched")
}

// Config parse options to config
func (opt *Options) Config() (*config.Config, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", opt.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("can't parse kubeconfig from (%v)", opt.Kubeconfig)
	}

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("can't new kube client: %v", err)
	}

	extClient, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("can't new extension client: %v", err)
	}

	var (
		kubeInformerOpts []informers.SharedInformerOption
		podInformerOpts  []informers.SharedInformerOption
		extInformerOpts  []extinformers.SharedInformerOption
	)

	if len(opt.Namespace) != 0 {
		kubeInformerOpts = append(kubeInformerOpts, informers.WithNamespace(opt.Namespace))
		podInformerOpts = append(podInformerOpts, informers.WithNamespace(opt.Namespace))
		extInformerOpts = append(extInformerOpts, extinformers.WithNamespace(opt.Namespace))
	}

	// only watch changes of pods in mario stage
	podInformerOpts = append(podInformerOpts, informers.WithTweakListOptions(
		func(opts *metav1.ListOptions) {
			opts.LabelSelector =
				labels.SelectorFromValidatedSet(
					labels.Set{
						v1alpha1.DefaultFlowStageLabelKey: v1alpha1.FlowStageMario,
					},
				).String()
		},
	))

	kubeInformerFactory := informers.NewSharedInformerFactoryWithOptions(kubeClient, 0, kubeInformerOpts...)
	podInformerFactory := informers.NewSharedInformerFactoryWithOptions(kubeClient, 0, podInformerOpts...)
	extInformerFactory := extinformers.NewSharedInformerFactoryWithOptions(extClient, 0, extInformerOpts...)

	eventInformer := extInformerFactory.Mario().V1alpha1().Events()
	pipeInformer := extInformerFactory.Mario().V1alpha1().Pipes()
	flowInformer := extInformerFactory.Mario().V1alpha1().Flows()

	jobInformer := kubeInformerFactory.Batch().V1().Jobs()
	pvcInformer := kubeInformerFactory.Core().V1().PersistentVolumeClaims()
	cmInformer := kubeInformerFactory.Core().V1().ConfigMaps()
	podInformer := podInformerFactory.Core().V1().Pods()

	c := &config.Config{
		KubeClient: kubeClient,
		ExtClient:  extClient,

		KubeInformerFactory: kubeInformerFactory,
		PodInformerFactory:  podInformerFactory,
		ExtInformerFactory:  extInformerFactory,

		EventInformer: eventInformer,
		PipeInformer:  pipeInformer,
		FlowInformer:  flowInformer,

		JobInformer:       jobInformer,
		PVCInformer:       pvcInformer,
		ConfigMapInformer: cmInformer,
		PodInformer:       podInformer,
	}

	return c, nil
}
