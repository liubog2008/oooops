package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/liubog2008/pkg/http/errors"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/hook"
	v3 "github.com/liubog2008/oooops/pkg/hook/github/api/v3"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

var (
	ErrMethodNotAllowed = errors.MustNewFactory(http.StatusMethodNotAllowed, "MethodNotAllowed", "method not allowed")

	ErrMissingEventType = errors.MustNewFactory(http.StatusBadRequest, "MissingEventType", "missing header "+v3.HeaderGithubEvent)

	ErrMissingDelivery = errors.MustNewFactory(http.StatusBadRequest, "MissingDelivery", "missing header "+v3.HeaderGithubDelivery)

	ErrMissingSignature = errors.MustNewFactory(http.StatusForbidden, "MissingSignature", "missing header "+v3.HeaderHubSignature)

	ErrInvalidSignature = errors.MustNewFactory(http.StatusForbidden, "InvalidSignature", v3.HeaderHubSignature+" is invalid")

	ErrUnsupportedContentType = errors.MustNewFactory(
		http.StatusBadRequest,
		"UnsupportedContentType",
		"only ["+ContentTypeJSON+"] is supported",
	)

	ErrFailToUnmarshal = errors.MustNewFactory(http.StatusBadRequest, "FailToUnmarshal", "failed to unmarshal body: %{err}")
)

type githubHook struct {
	token string
	wg    sync.WaitGroup
}

func New(token string) hook.Hook {
	return &githubHook{
		token: token,
	}
}

func (h *githubHook) Dispatch(req *http.Request) error {
	eventType, eventGUID, payload, err := h.validate(req)
	if err != nil {
		return err
	}

	switch eventType {
	case v3.EventTypeIssueComment:
		var e v3.IssueCommentEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return ErrFailToUnmarshal.New(err)
		}
		e.GUID = eventGUID
		h.wg.Add(1)
		go h.handleIssueCommentEvent(&e)
	case v3.EventTypePullRequest:
		var e v3.PullRequestEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return ErrFailToUnmarshal.New(err)
		}
		e.GUID = eventGUID
		h.wg.Add(1)
		go h.handlePullRequestEvent(&e)
	case v3.EventTypeCreate:
		var e v3.CreateEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return ErrFailToUnmarshal.New(err)
		}
		e.GUID = eventGUID
		h.wg.Add(1)
		go h.handleCreateEvent(&e)
	default:
		klog.Infof("[%s] unsupported event type %s", eventGUID, eventType)
	}
	return nil
}

func (h *githubHook) validate(req *http.Request) (v3.EventType, string, []byte, error) {
	defer req.Body.Close()

	if req.Method != http.MethodPost {
		return "", "", nil, ErrMethodNotAllowed.New()
	}
	eventType := req.Header.Get(v3.HeaderGithubEvent)

	if eventType == "" {
		return "", "", nil, ErrMissingEventType.New()
	}

	eventGUID := req.Header.Get(v3.HeaderGithubDelivery)
	if eventGUID == "" {
		return "", "", nil, ErrMissingDelivery.New()
	}
	sig := req.Header.Get(v3.HeaderHubSignature)
	if sig == "" {
		return "", "", nil, ErrMissingSignature.New()
	}

	// token := req.Header.Get("X-Gitlab-Token")
	contentType := req.Header.Get(HeaderContentType)

	if contentType != ContentTypeJSON {
		return "", "", nil, ErrUnsupportedContentType.New()
	}
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", "", nil, err
	}

	if !validatePayload(payload, sig, h.token) {
		return "", "", nil, ErrInvalidSignature.New()
	}

	return v3.EventType(eventType), eventGUID, payload, nil
}

func validatePayload(payload []byte, sig string, token string) bool {
	if !strings.HasPrefix(sig, "sha1=") {
		return false
	}
	sig = sig[5:]
	decoded, err := hex.DecodeString(sig)
	if err != nil {
		return false
	}
	mac := hmac.New(sha1.New, []byte(token))
	mac.Write(payload)
	expected := mac.Sum(nil)
	return hmac.Equal(decoded, expected)
}
