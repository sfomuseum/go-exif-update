package rewrite

import (
	"bufio"
	"bytes"
	"golang.org/x/net/html"
	"io"
	_ "log"
	go_http "net/http"
	go_httptest "net/http/httptest"
	"strconv"
	"strings"
)

type RewriteHTMLFunc func(node *html.Node, writer io.Writer)

func RewriteHTMLHandler(prev go_http.Handler, rewrite_func RewriteHTMLFunc) go_http.Handler {

	fn := func(rsp go_http.ResponseWriter, req *go_http.Request) {

		rec := go_httptest.NewRecorder()
		prev.ServeHTTP(rec, req)

		prev_rsp := rec.Result()
		prev_headers := prev_rsp.Header

		defer prev_rsp.Body.Close()

		location := prev_headers.Get("Location")

		if location != "" {
			
			for k, v := range prev_headers {

				if k == "Location" {
					continue
				}
				
				for _, vv := range v {
					rsp.Header().Set(k, vv)
				}
			}
			
			go_http.Redirect(rsp, req, location, 303)
			return
		}
		
		content_type := prev_headers.Get("Content-Type")

		if content_type != "" {

			parts := strings.Split(content_type, ";")

			if parts[0] != "text/html" {

				for k, v := range prev_headers {

					for _, vv := range v {
						rsp.Header().Set(k, vv)
					}
				}

				_, err := io.Copy(rsp, prev_rsp.Body)

				if err != nil {
					go_http.Error(rsp, err.Error(), go_http.StatusInternalServerError)
					return
				}

				return
			}
		}

		doc, err := html.Parse(prev_rsp.Body)

		if err != nil {
			go_http.Error(rsp, err.Error(), go_http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		rewrite_func(doc, wr)

		err = html.Render(wr, doc)

		if err != nil {
			go_http.Error(rsp, err.Error(), go_http.StatusInternalServerError)
			return
		}

		wr.Flush()

		for k, v := range rec.Header() {

			if k == "Content-Length" {
				continue
			}

			if k == "Content-Type" {
				continue
			}

			rsp.Header()[k] = v
		}

		data := buf.Bytes()
		clen := len(data)

		rsp.Header().Set("Content-Length", strconv.Itoa(clen))
		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")

		rsp.WriteHeader(200)
		rsp.Write(data)
	}

	return go_http.HandlerFunc(fn)
}
