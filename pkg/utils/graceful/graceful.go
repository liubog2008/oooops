// Package graceful defines helper functions to let program gracefully shutdown
package graceful

import (
	"context"
	"sync"
	"time"

	"k8s.io/klog"
)

// Interface defines interface to help graceful shutdown program
type Interface interface {
	OnShutdown(func(context.Context))
	WaitForShutdown(<-chan struct{}, time.Duration)
}

type graceful struct {
	wg      sync.WaitGroup
	workers []func(context.Context)
}

func New() Interface {
	return &graceful{}
}

// OnShutdown register worker which will be called after stopCh is closed
func (g *graceful) OnShutdown(worker func(context.Context)) {
	g.wg.Add(1)
	g.workers = append(g.workers, worker)
}

// WaitForShutdown will wait until all workers are done or timeout
func (g *graceful) WaitForShutdown(stopCh <-chan struct{}, gracePeriod time.Duration) {
	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	for i := range g.workers {
		go func(k int) {
			w := g.workers[k]
			w(ctx)
			klog.Infof("worker %d is done", k)
			g.wg.Done()
		}(i)
	}
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		klog.Infof("all workers are gracefully terminated")
	case <-ctx.Done():
		klog.Infof("timeout to wait for graceful terminating")
	}
}
