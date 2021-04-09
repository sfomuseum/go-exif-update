package rewrite

import (
	go_http "net/http"
)

type RewriteRequestFunc func(req *go_http.Request) (*go_http.Request, error)

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
