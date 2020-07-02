package pipe

import (
	"fmt"
	"reflect"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) addEvent(obj interface{}) {
	event, ok := obj.(*v1alpha1.Event)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("obj is not Event: %v", obj))
		return
	}
	pipes, err := c.getPipeWatchers(event)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("can't get pipe watchers of event %s/%s: %v", event.Namespace, event.Name, err))
		return
	}

	for _, pipe := range pipes {
		c.addPipe(pipe)
	}
}

func (c *Controller) updateEvent(old, cur interface{}) {
	oldEvent, ok1 := old.(*v1alpha1.Event)
	curEvent, ok2 := cur.(*v1alpha1.Event)
	if !ok1 || !ok2 {
		utilruntime.HandleError(fmt.Errorf("either old or cur is not Event: %v, %v", old, cur))
		return
	}
	if oldEvent.ResourceVersion == curEvent.ResourceVersion {
		return
	}
	if reflect.DeepEqual(&oldEvent.Spec, &curEvent.Spec) {
		return
	}
	c.addEvent(curEvent)
}

func (c *Controller) deleteEvent(obj interface{}) {
	if event, ok := obj.(*v1alpha1.Event); ok {
		c.addEvent(event)
		return
	}
	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
		return
	}
	event, ok := tombstone.Obj.(*v1alpha1.Event)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not an Event: %#v", tombstone.Obj))
		return
	}
	c.addEvent(event)
}

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
		c.addPipe(pipe)
		return
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
		return
	}
	pipe, ok := tombstone.Obj.(*v1alpha1.Pipe)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Pipe: %#v", tombstone.Obj))
		return
	}
	c.addPipe(pipe)
}
