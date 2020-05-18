// Package flow defines a flow controller to generate job and sync flow status
package flow

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	batchlisters "k8s.io/client-go/listers/batch/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/client/clientset"
	"github.com/liubog2008/oooops/pkg/client/clientset/scheme"
	marioinformers "github.com/liubog2008/oooops/pkg/client/informers/mario/v1alpha1"
	mariolisters "github.com/liubog2008/oooops/pkg/client/listers/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/controller"
)

const (
	gitRootVolumeName = "git"

	marioWorkingDir = "/repo"
)

// ControllerOptions defines options which is needed by flow controller
type ControllerOptions struct {
	KubeClient kubernetes.Interface

	ExtClient clientset.Interface

	JobInformer batchinformers.JobInformer

	FlowInformer marioinformers.FlowInformer

	PVCInformer coreinformers.PersistentVolumeClaimInformer
}

// Controller defines controller to manage flow lifecycle and generate jobs
// It will do these things:
// - Generate jobs of the flow
// - Calculate flow status
type Controller struct {
	schema.GroupVersionKind

	kubeClient kubernetes.Interface
	extClient  clientset.Interface

	flowLister mariolisters.FlowLister
	jobLister  batchlisters.JobLister
	pvcLister  corelisters.PersistentVolumeClaimLister

	informersSynced []cache.InformerSynced

	eventBroadcaster record.EventBroadcaster
	eventRecorder    record.EventRecorder

	queue workqueue.RateLimitingInterface

	buildReconciler controller.ReconcilerBuilder

	marioImage string
	gitImage   string
	gitCommand []string
}

// NewController returns a flow controller
func NewController(opt *ControllerOptions) *Controller {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartLogging(klog.Infof)
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: opt.KubeClient.CoreV1().Events("")})
	recorder := broadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "flow"})

	c := &Controller{
		GroupVersionKind: v1alpha1.SchemeGroupVersion.WithKind("Flow"),

		kubeClient: opt.KubeClient,
		extClient:  opt.ExtClient,

		informersSynced: []cache.InformerSynced{
			opt.JobInformer.Informer().HasSynced,
			opt.FlowInformer.Informer().HasSynced,
			opt.PVCInformer.Informer().HasSynced,
		},

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "flow"),

		flowLister: opt.FlowInformer.Lister(),
		jobLister:  opt.JobInformer.Lister(),
		pvcLister:  opt.PVCInformer.Lister(),

		eventBroadcaster: broadcaster,
		eventRecorder:    recorder,

		buildReconciler: controller.BuildRateLimitingReconciler,

		gitImage:   "alpine/git:v2.24.3",
		marioImage: "busybox",
	}

	opt.FlowInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addFlow,
		UpdateFunc: c.updateFlow,
		DeleteFunc: c.deleteFlow,
	})

	opt.JobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addJob,
		UpdateFunc: c.updateJob,
		DeleteFunc: c.deleteJob,
	})

	return c
}

// Run will start the controller
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting flow controller")
	defer klog.Infof("Shutting down flow controller")

	if !cache.WaitForCacheSync(stopCh, c.informersSynced...) {
		utilruntime.HandleError(fmt.Errorf("unable to sync caches for flow controller"))
		return
	}

	klog.Infof("Cache of flow controller has been synced")

	for i := 0; i < workers; i++ {
		controller.WaitUntil("flow", c.buildReconciler(c.queue, c.syncFlow), stopCh)
	}

	klog.Infof("flow controller is working")

	<-stopCh
}
