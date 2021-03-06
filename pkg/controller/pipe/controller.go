package pipe

import (
	"fmt"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/client/clientset"
	"github.com/liubog2008/oooops/pkg/client/clientset/scheme"
	marioinformers "github.com/liubog2008/oooops/pkg/client/informers/mario/v1alpha1"
	mariolisters "github.com/liubog2008/oooops/pkg/client/listers/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

// ControllerOptions defines options of the controller
type ControllerOptions struct {
	KubeClient kubernetes.Interface

	ExtClient clientset.Interface

	EventInformer marioinformers.EventInformer

	PipeInformer marioinformers.PipeInformer

	FlowInformer marioinformers.FlowInformer
}

// Controller defines controller to manage pipe lifecycle and generate flow
type Controller struct {
	schema.GroupVersionKind

	kubeClient kubernetes.Interface
	extClient  clientset.Interface

	pipeLister  mariolisters.PipeLister
	eventLister mariolisters.EventLister
	flowLister  mariolisters.FlowLister

	informersSynced []cache.InformerSynced

	eventBroadcaster record.EventBroadcaster
	eventRecorder    record.EventRecorder

	queue workqueue.RateLimitingInterface

	buildReconciler controller.ReconcilerBuilder
}

// NewController returns the pipe controller
func NewController(opt *ControllerOptions) *Controller {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartLogging(klog.Infof)
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: opt.KubeClient.CoreV1().Events("")})
	recorder := broadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "pipe"})

	c := &Controller{
		GroupVersionKind: v1alpha1.SchemeGroupVersion.WithKind("Pipe"),

		kubeClient: opt.KubeClient,
		extClient:  opt.ExtClient,

		informersSynced: []cache.InformerSynced{
			opt.PipeInformer.Informer().HasSynced,
			opt.EventInformer.Informer().HasSynced,
			opt.FlowInformer.Informer().HasSynced,
		},

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pipe"),

		pipeLister:  opt.PipeInformer.Lister(),
		eventLister: opt.EventInformer.Lister(),
		flowLister:  opt.FlowInformer.Lister(),

		eventBroadcaster: broadcaster,
		eventRecorder:    recorder,

		buildReconciler: controller.BuildRateLimitingReconciler,
	}

	opt.PipeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addPipe,
		UpdateFunc: c.updatePipe,
		DeleteFunc: c.deletePipe,
	})

	opt.EventInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addEvent,
		UpdateFunc: c.updateEvent,
		DeleteFunc: c.deleteEvent,
	})

	return c
}

// Run will start the controller
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting pipe controller")
	defer klog.Infof("Shutting down pipe controller")

	if !cache.WaitForCacheSync(stopCh, c.informersSynced...) {
		utilruntime.HandleError(fmt.Errorf("unable to sync caches for pipe controller"))
		return
	}

	klog.Infof("Cache of pipe controller has been synced")

	for i := 0; i < workers; i++ {
		controller.WaitUntil("pipe", c.buildReconciler(c.queue, c.syncPipe), stopCh)
	}

	klog.Infof("pipe controller is working")

	<-stopCh
}
