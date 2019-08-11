package flow

import (
	"fmt"
	"testing"
	"time"

	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	"github.com/liubog2008/oooops/pkg/client/clientset"
	"github.com/liubog2008/oooops/pkg/client/clientset/fake"
	"github.com/liubog2008/oooops/pkg/client/informers"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubefake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/util/workqueue"
)

type fakeController struct {
	fake         *Controller
	syncStarter  chan struct{}
	syncFinished chan struct{}
}

func newFakeController(kubeClient kubernetes.Interface, extClient clientset.Interface) *fakeController {
	fakeExtInformers := informers.NewSharedInformerFactory(f.extClient, noResyncPeriodFunc())
	fakeKubeInformers := kubeinformers.NewSharedInformerFactory(f.kubeClient, noResyncPeriodFunc())

	fakeFlowInformer := fakeExtInformers.Flow().V1alpha1().Flows()
	fakeJobInformer := fakeKubeInformers.Batch().V1().Jobs()

	return &fakeController{
		fake: NewController(&ControllerOptions{
			KubeClient:   kubeClient,
			ExtClient:    extClient,
			FlowInformer: fakeFlowInformer,
			JobInformer:  fakeJobInformer,
		}),
		syncStarter:  make(chan struct{}),
		syncFinished: make(chan struct{}),
	}
}

func (c *fakeController) syncOnce(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	c.syncStarter <- struct{}{}
	select {
	case <-c.syncFinished:
	case <-timer.C:
		return fmt.Errorf("sync timeout")
	}
	return nil
}

func (c *fakeController) MustSyncOnce(t *testing.T) {
	err := c.syncOnce(defaultTimeout)
	require.NoError(t, err, "sync once timeout")
}

func (c *fakeController) worker(queue workqueue.RateLimitingInterface, handler func(key string) error) func() bool {
	w := c.fake.workerFactory(queue, handler)
	return func() bool {
		<-c.syncStarter
		quit := w()
		c.syncFinished <- struct{}{}
		return quit
	}
}

func (c *fakeController) Run(stopCh <-chan struct{}) {
	w := c.fake.workerFactory
	c.fake.workerFactory = c.worker
	defer func() {
		c.fake.workerFactory = w
	}()

	c.fake.Run(1, stopCh)
}

const defaultTimeout = time.Second * 60

func TestController(t *testing.T) {
	stopCh := make(chan struct{})
	kubeStore := []runtime.Object{}
	extStore := []runtime.Object{}
	kubeClient := kubefake.NewSimpleClientset(kubeStore...)
	extClient := fake.NewSimpleClientset(extStore...)

	fc := newFakeController(kubeClient, extClient)
	go fc.Run(stopCh)

	_, err := extClient.FlowV1alpha1().Flows().Create(&v1alpha1.Flow{})
	require.NoError(t, err, "create flow error")
	fc.MustSyncOnce(t)

	kubeClient.Actions()

	utiltesting.AssertAction(t, expected, actions[len(creations):])

	stopCh <- struct{}{}
}
