package mario

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/git"
	"github.com/liubog2008/pkg/http/errors"
	"k8s.io/klog"
)

const (
	authKey = "Authorization"

	tokenType = "Bearer"
)

var (
	// ErrUnauthorized defines error which means token is not right
	ErrUnauthorized = errors.MustNewFactory(http.StatusUnauthorized, "Unauthorized", "unauthorized: %{err}")

	// ErrHasBeenConsumed defines error that mario file has been consumed by others
	ErrHasBeenConsumed = errors.MustNewFactory(http.StatusUnprocessableEntity, "HasBeenConumsed", "mario file has been consumed")
)

type mario struct {
	gitCmd git.Interface

	workDir string

	token string

	gracefulShutdownTimeout time.Duration

	lock sync.Mutex

	counter int

	addr string

	// wait is a channel used to wait after net listen
	// it is used for testing
	wait chan<- struct{}
}

func newWithWait(local, remote, addr, token string, wait chan<- struct{}) (Interface, error) {
	gitCmd, err := git.New(local)
	if err != nil {
		return nil, err
	}

	if err := gitCmd.WithRepo(remote); err != nil {
		return nil, err
	}

	m := mario{
		gitCmd:  gitCmd,
		workDir: local,
		addr:    addr,
		token:   token,
		wait:    wait,
	}

	return &m, nil
}

// New returns a mario interface
func New(local, remote, addr, token string) (Interface, error) {
	return newWithWait(local, remote, addr, token, nil)
}

func (m *mario) Checkout(ref string) error {
	if err := m.gitCmd.Fetch(ref); err != nil {
		return err
	}

	if err := m.gitCmd.Clean(); err != nil {
		return err
	}

	return m.gitCmd.Checkout(ref)
}

func (m *mario) Serve(stopCh chan struct{}) error {
	klog.Infof("mario begin to serve file")
	startTime := time.Now()
	defer func() {
		klog.Infof("mario serve file finished, cost (%v)", time.Since(startTime))
	}()

	router := http.NewServeMux()
	router.HandleFunc("/", m.handleFunc(stopCh))

	svr := &http.Server{
		Addr:         m.addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan struct{})

	ln, err := net.Listen("tcp", m.addr)
	if err != nil {
		return err
	}

	if m.wait != nil {
		close(m.wait)
	}

	go waitToShutdown(svr, m.gracefulShutdownTimeout, stopCh, done)

	if err := svr.Serve(ln); err != nil && err != http.ErrServerClosed {
		klog.Errorf("Could not listen on %s: %v", m.addr, err)
		return err
	}

	<-done

	return nil
}

func waitToShutdown(server *http.Server, timeout time.Duration, stopCh <-chan struct{}, done chan<- struct{}) {
	<-stopCh
	klog.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		klog.Fatalf("Could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

func (m *mario) handleFunc(stopCh chan<- struct{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := m.checkToken(r); err != nil {
			writeError(w, ErrUnauthorized.New(err))
			return
		}

		m.lock.Lock()

		if m.counter != 0 {
			writeError(w, ErrHasBeenConsumed.New())
			m.lock.Unlock()
			return
		}

		m.counter++
		m.lock.Unlock()

		http.ServeFile(w, r, filepath.Join(m.workDir, v1alpha1.MarioFile))
		close(stopCh)
	}
}

func (m *mario) checkToken(r *http.Request) error {
	v := r.Header.Get(authKey)
	if v == "" {
		return fmt.Errorf("token is not found, please set token into Authorization header")
	}
	typeAndToken := strings.Split(v, " ")
	if len(typeAndToken) != 2 {
		return fmt.Errorf("bad token format, expected: Bearer $token, actual: %s", v)
	}
	typ := strings.TrimSpace(typeAndToken[0])
	if typ != tokenType {
		return fmt.Errorf("bad token format, invalid token type, expected: %s, actual: %s", tokenType, typ)
	}
	if strings.TrimSpace(typeAndToken[1]) != m.token {
		return fmt.Errorf("token is not equal")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	e := unwarp(err)
	b, err := json.Marshal(e)
	if err != nil {
		klog.Errorf("can't marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(e.Code)
	if _, err := w.Write(b); err != nil {
		klog.Errorf("can't write whole response: %v", err)
	}
}

func unwarp(err error) *errors.Error {
	switch e := err.(type) {
	case *errors.Error:
		return e
	default:
		return &errors.Error{
			Code:    http.StatusInternalServerError,
			Reason:  "Unknown",
			Message: err.Error(),
		}
	}
}
