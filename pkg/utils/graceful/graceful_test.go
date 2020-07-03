package graceful

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGracefulShutdown(t *testing.T) {
	stopCh := make(chan struct{})
	done := make(chan struct{})
	var count int32
	workerNum := 10

	go func() {
		g := New()

		for i := 0; i < workerNum; i++ {
			g.OnShutdown(func(ctx context.Context) {
				atomic.AddInt32(&count, 1)
			})
		}

		g.WaitForShutdown(stopCh, 100*time.Millisecond)

		close(done)
	}()

	close(stopCh)

	<-done

	assert.Equal(t, int(count), workerNum)
}

func TestGracefulShutdownTimeout(t *testing.T) {
	stopCh := make(chan struct{})
	done := make(chan struct{})
	var count int32
	workerNum := 10

	go func() {
		g := New()

		for i := 0; i < workerNum; i++ {
			g.OnShutdown(func(ctx context.Context) {
				// sleep 2 sec to ensure that shutdown will timeout
				time.Sleep(time.Second)
				atomic.AddInt32(&count, 1)
			})
		}

		g.WaitForShutdown(stopCh, 100*time.Millisecond)

		close(done)
	}()

	close(stopCh)

	<-done

	assert.Equal(t, int(count), 0)
}
