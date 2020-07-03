package mario

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/liubog2008/pkg/http/errors"
	"github.com/munnerz/goautoneg"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/client/clientset/scheme"
	"github.com/liubog2008/oooops/pkg/mario/git"
	"github.com/liubog2008/oooops/pkg/utils/graceful"
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

	// ErrNotAcceptable defines error that server can't handle current content type
	ErrNotAcceptable = errors.MustNewFactory(
		http.StatusNotAcceptable,
		"NotAcceptable",
		"only these media types [%{supported}] are supported, but client only accept [%{accept}]",
	)

	// ErrEncoding defines error that server fail to encode mario to specified content type
	ErrEncoding = errors.MustNewFactory(
		http.StatusInternalServerError,
		"FailedToEncode",
		"can't encode mario to %{contentType}: %{err}",
	)
)

// Interface defines interface to verify and serve mario file
type Interface interface {
	Run(stopCh <-chan struct{}) error
}

// Config defines config to run mario server
type Config struct {
	GitCommand              git.Interface
	Addr                    string
	GracefulShutdownTimeout time.Duration

	Remote string
	Ref    string

	Token string
}

type mario struct {
	gitCmd                  git.Interface
	addr                    string
	gracefulShutdownTimeout time.Duration
	remote                  string
	ref                     string

	token string

	lock    sync.Mutex
	counter int

	obj *v1alpha1.Mario
}

// New returns a mario interface
func New(c *Config) Interface {
	m := mario{
		gitCmd:                  c.GitCommand,
		addr:                    c.Addr,
		gracefulShutdownTimeout: c.GracefulShutdownTimeout,
		remote:                  c.Remote,
		ref:                     c.Ref,
		token:                   c.Token,
	}

	return &m
}

func (m *mario) Run(stopCh <-chan struct{}) error {
	if err := m.gitCmd.Verify(m.remote, m.ref); err != nil {
		return err
	}

	body, err := ioutil.ReadFile(v1alpha1.MarioFile)
	if err != nil {
		return err
	}

	decoder := scheme.Codecs.UniversalDecoder(v1alpha1.SchemeGroupVersion)

	marioObj := v1alpha1.Mario{}

	if _, _, err := decoder.Decode(body, nil, &marioObj); err != nil {
		return err
	}

	m.obj = &marioObj

	return m.serve(stopCh)
}

func (m *mario) serve(stopCh <-chan struct{}) error {
	klog.Infof("mario begin to serve file")
	startTime := time.Now()
	defer func() {
		klog.Infof("mario serve file finished, cost (%v)", time.Since(startTime))
	}()

	done := make(chan struct{})

	go func() {
		<-stopCh
		close(done)
	}()

	router := http.NewServeMux()
	router.HandleFunc("/healthz", m.health)
	router.HandleFunc("/", m.handleFunc(done))

	srv := &http.Server{
		Addr:         m.addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	g := graceful.New()

	defer g.WaitForShutdown(done, m.gracefulShutdownTimeout)

	g.OnShutdown(func(ctx context.Context) {
		if err := srv.Shutdown(ctx); err != nil {
			klog.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			klog.Infof("listen and serve finished: %v", err)
		}
	}()

	return nil
}

func (m *mario) handleFunc(done chan<- struct{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := m.checkToken(r); err != nil {
			writeError(w, ErrUnauthorized.New(err))
			return
		}
		supported := scheme.Codecs.SupportedMediaTypes()
		accept := r.Header.Get("Accept")
		info := isAcceptable(accept, supported)
		if info == nil {
			mediaTypes := []string{}
			for i := range supported {
				info := &supported[i]
				mediaTypes = append(mediaTypes, info.MediaType)
			}
			writeError(w, ErrNotAcceptable.New(mediaTypes, accept))
			return
		}
		s := info.Serializer
		encoder := scheme.Codecs.EncoderForVersion(s, v1alpha1.SchemeGroupVersion)

		m.lock.Lock()

		if m.counter != 0 {
			writeError(w, ErrHasBeenConsumed.New())
			m.lock.Unlock()
			return
		}

		m.counter++
		m.lock.Unlock()

		if err := encoder.Encode(m.obj, w); err != nil {
			writeError(w, ErrEncoding.New(info.MediaType, err))
		}

		close(done)
	}
}

func isAcceptable(header string, accepted []runtime.SerializerInfo) *runtime.SerializerInfo {
	if len(header) == 0 && len(accepted) > 0 {
		return &accepted[0]
	}

	clauses := goautoneg.ParseAccept(header)
	for i := range clauses {
		clause := &clauses[i]
		for i := range accepted {
			accepts := &accepted[i]
			switch {
			case clause.Type == accepts.MediaTypeType && clause.SubType == accepts.MediaTypeSubType,
				clause.Type == accepts.MediaTypeType && clause.SubType == "*",
				clause.Type == "*" && clause.SubType == "*":
				return accepts
			}
		}
	}

	return nil
}

func (m *mario) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		klog.Errorf("can't write response: %v", err)
	}
}

func (m *mario) checkToken(r *http.Request) error {
	v := r.Header.Get(authKey)
	if v == "" {
		return fmt.Errorf("token is not found, please set token into Authorization header")
	}
	typeAndToken := strings.SplitN(v, " ", 2)
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
