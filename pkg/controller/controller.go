package controller

import (
	"time"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

// Syncer is the main function to sync desired state and actual state of the key
type Syncer func(key string) error

// Reconciler defines function to sync and handle sync result
type Reconciler func() (bool, error)

// ReconcilerBuilder defines factory function returns Reconciler
type ReconcilerBuilder func(queue workqueue.RateLimitingInterface, syncer Syncer) Reconciler

// BuildRateLimitingReconciler defines factory to return rate limiting reconciler
func BuildRateLimitingReconciler(queue workqueue.RateLimitingInterface, syncer Syncer) Reconciler {
	return func() (bool, error) {
		key, quit := queue.Get()
		if quit {
			return true, nil
		}
		defer queue.Done(key)

		if err := syncer(key.(string)); err != nil {
			queue.AddRateLimited(key)
			return false, err
		}
		queue.Forget(key)
		return false, nil
	}
}

// WaitUntil defines a main loop of controller
func WaitUntil(name string, reconciler Reconciler, stopCh <-chan struct{}) {
	forever := func() {
		for {
			quit, err := reconciler()
			if err != nil {
				utilruntime.HandleError(err)
			}

			if quit {
				klog.Infof("%s controller worker shutting down", name)
				return
			}
		}
	}
	go wait.Until(forever, time.Second, stopCh)
}
