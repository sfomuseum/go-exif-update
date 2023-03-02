package rewrite

import (
	go_http "net/http"
)

// RewriteRequestFunc is a custom callback function for creating a new `http.Request` instance derived from 'req'.
type RewriteRequestFunc func(req *go_http.Request) (*go_http.Request, error)

// RewriteRequestHandler() creates a `net/http` middleware handler that invokes 'rewrite_func' with
// the current request to create a new or updated `http.Request` instance used to serve 'next'.
func RewriteRequestHandler(next go_http.Handler, rewrite_func RewriteRequestFunc) go_http.Handler {

	fn := func(rsp go_http.ResponseWriter, req *go_http.Request) {

		next_req, err := rewrite_func(req)

		if err != nil {
			go_http.Error(rsp, err.Error(), go_http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(rsp, next_req)
	}

	return go_http.HandlerFunc(fn)
}
