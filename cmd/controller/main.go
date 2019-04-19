package main

import (
	"time"

	"github.com/liubog2008/oooops/pkg/client/clientset"
	flowinformers "github.com/liubog2008/oooops/pkg/client/informers"
	"github.com/liubog2008/oooops/pkg/controller/flow"
	"github.com/spf13/pflag"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var kubeconfig string

func init() {
	klog.InitFlags(nil)
}

func main() {
	pflag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	pflag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	// create the clientset
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	extclient, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	stopCh := make(chan struct{})

	kubeSharedInformerFactory := informers.NewSharedInformerFactory(kubeclient, 5*time.Minute)
	extSharedInformerFactory := flowinformers.NewSharedInformerFactory(extclient, 5*time.Minute)
	jobInformer := kubeSharedInformerFactory.Batch().V1().Jobs()
	flowInformer := extSharedInformerFactory.Flow().V1alpha1().Flows()

	c := flow.NewController(&flow.ControllerOptions{
		KubeClient:   kubeclient,
		ExtClient:    extclient,
		JobInformer:  jobInformer,
		FlowInformer: flowInformer,
	})

	klog.Infof("start shared informer factory")
	kubeSharedInformerFactory.Start(stopCh)
	extSharedInformerFactory.Start(stopCh)
	c.Run(1, stopCh)
	close(stopCh)
}
