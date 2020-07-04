package hook

import "net/http"

type Hook interface {
	Dispatch(req *http.Request) error
}
