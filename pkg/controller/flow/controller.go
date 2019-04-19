package flow

import (
	"fmt"
	"time"

	"github.com/liubog2008/oooops/pkg/client/clientset"
	flowinformers "github.com/liubog2008/oooops/pkg/client/informers/flow/v1alpha1"
	flowlisters "github.com/liubog2008/oooops/pkg/client/listers/flow/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/kubernetes"
	batchlisters "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

// ControllerOptions defines options to flow controller
type ControllerOptions struct {
	// KubeClient defines interface of raw kubernetes API
	KubeClient kubernetes.Interface

	// ExtClient defines interface of extension API of this project
	ExtClient clientset.Interface

	// FlowInformer defines informer of flow
	FlowInformer flowinformers.FlowInformer

	// JobInformer defines informer of job
	JobInformer batchinformers.JobInformer

	// SCMImage defines image for scm
	SCMImage string
}

// Controller defines flow controller to control CI/CD flow
type Controller struct {
	kubeclient kubernetes.Interface
	extclient  clientset.Interface

	informersSynced []cache.InformerSynced

	flowQueue  workqueue.RateLimitingInterface
	flowLister flowlisters.FlowLister

	jobQueue  workqueue.RateLimitingInterface
	jobLister batchlisters.JobLister

	scmImage string

	jobCaches map[string]JobCache
}

// JobCache defines cache of jobs
type JobCache struct {
	hash  string
	cache map[string]*batchv1.Job
}

// NewController return a flow controller
func NewController(options *ControllerOptions) *Controller {
	c := Controller{
		kubeclient: options.KubeClient,
		extclient:  options.ExtClient,
		flowLister: options.FlowInformer.Lister(),
		jobLister:  options.JobInformer.Lister(),
		informersSynced: []cache.InformerSynced{
			options.FlowInformer.Informer().HasSynced,
			options.JobInformer.Informer().HasSynced,
		},
		flowQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "flow"),
		jobQueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "flow job"),

		jobCaches: map[string]JobCache{},
		scmImage:  options.SCMImage,
	}

	options.FlowInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addFlow,
			UpdateFunc: c.updateFlow,
			// DeleteFunc: c.deleteFlow,
		},
	)

	options.JobInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addJob,
			UpdateFunc: c.updateJob,
			DeleteFunc: c.deleteJob,
		},
	)
	return &c
}

func (c *Controller) worker(queue workqueue.RateLimitingInterface, handler func(key string) error) func() {
	workFunc := func() bool {
		key, quit := queue.Get()
		if quit {
			return true
		}
		defer queue.Done(key)

		if err := handler(key.(string)); err != nil {
			utilruntime.HandleError(err)
			queue.AddRateLimited(key)
			return false
		}
		queue.Forget(key)
		return false
	}

	return func() {
		for {
			if quit := workFunc(); quit {
				klog.Infof("flow controller worker shutting down")
				return
			}
		}
	}
}

// Run will start the flow controller
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.flowQueue.ShutDown()

	klog.Infof("Starting flow controller")
	defer klog.Infof("Shutting down flow controller")

	if !cache.WaitForCacheSync(stopCh, c.informersSynced...) {
		utilruntime.HandleError(fmt.Errorf("Unable to sync caches for flow controller"))
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker(c.flowQueue, c.syncFlowHandler), time.Second, stopCh)
	}
	<-stopCh
}
