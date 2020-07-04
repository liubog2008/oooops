package hook

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/liubog2008/pkg/http/errors"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/hook"
	"github.com/liubog2008/oooops/pkg/hook/github"
	"github.com/liubog2008/oooops/pkg/utils/graceful"
)

type handler struct {
	h hook.Hook
}

func writeError(w http.ResponseWriter, err error) {
	e := errors.Assert(err)
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

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.h.Dispatch(r); err != nil {
		writeError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	h := handler{
		h: github.New("test"),
	}

	router := http.NewServeMux()
	router.Handle("/hook", &h)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	stopCh := make(chan struct{})

	g := graceful.New()
	defer g.WaitForShutdown(stopCh, 30*time.Second)

	g.OnShutdown(func(ctx context.Context) {
		if err := srv.Shutdown(ctx); err != nil {
			klog.Errorf("can't shutdown gracefully: %v", err)
		}
	})

	klog.Info(srv.ListenAndServe())
}
