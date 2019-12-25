package pipe

import (
	"fmt"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) addPipe(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}

	c.queue.Add(key)
}

func (c *Controller) updatePipe(old, cur interface{}) {
	c.addPipe(cur)
}

func (c *Controller) deletePipe(obj interface{}) {
	if pipe, ok := obj.(*v1alpha1.Pipe); ok {
		// Enqueue all the services that the pod used to be a member of.
		// This is the same thing we do when we add a pod.
		c.addPipe(pipe)
	}
	// If we reached here it means the pod was deleted but its final state is unrecorded.
	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Couldn't get object from tombstone %#v", obj))
		return
	}
	pipe, ok := tombstone.Obj.(*v1alpha1.Pipe)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Tombstone contained object that is not a Pipe: %#v", obj))
		return
	}
	c.addPipe(pipe)
}
