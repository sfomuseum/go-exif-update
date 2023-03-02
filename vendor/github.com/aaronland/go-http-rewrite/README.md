# go-http-rewrite

Go package for creating `net/http` middleware handlers to alter the contents of other (`net/http`) handlers.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-http-rewrite.svg)](https://pkg.go.dev/github.com/aaronland/go-http-rewrite)

## Example

_Error handler omitted for brevity._

### AppendResourcesHandler

Middleware handler to append custom JavaScript and CSS tags or data attributes to the HTML output of a previous handler.

```
package main

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func baseAppendHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Header().Set("Content-type", "text/html")
		rsp.Write([]byte(`<html><head><title>Test</title></head><body>Hello world</body><html>`))
	}

	h := http.HandlerFunc(fn)
	return h
}

func main() {

	append_opts := &AppendResourcesOptions{
		JavaScript:     []string{"test.js"},
		Stylesheets:    []string{"test.css"},
		DataAttributes: map[string]string{"example": "example"},
	}

	append_handler := baseAppendHandler()

	append_handler = AppendResourcesHandler(append_handler, append_opts)

	s := &http.Server{
		Addr:    ":8080",
		Handler: append_handler,
	}

	defer s.Close()

	go s.ListenAndServe()

	rsp, _ := http.Get("http://localhost:8080")

	defer rsp.Body.Close()

	body, _ := io.ReadAll(rsp.Body)
	fmt.Println(string(body))

	// Prints:
	// <html><head><title>Test</title><script type="text/javascript" src="test.js"></script><link type="text/css" rel="stylesheet" href="test.css"/></head><body data-example="example">Hello world</body></html>
	
}
```

### RewriteHTMLHandler

Middleware handler to rewrite the HTML output of a previous handler.

```
package main

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"testing"
)

func baseRewriteHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Header().Set("Content-type", "text/html")
		rsp.Write([]byte(`<html><head><title>Test</title></head><body><p>hello world</p></body></html>`))
	}

	h := http.HandlerFunc(fn)
	return h
}

func main() {

	var rewrite_func RewriteHTMLFunc

	rewrite_func = func(n *html.Node, wr io.Writer) {

		if n.Type == html.ElementNode && n.Data == "p" {
			n.FirstChild.Data = "HELLO WORLD"
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rewrite_func(c, wr)
		}

	}

	rewrite_handler := baseRewriteHandler()

	rewrite_handler = RewriteHTMLHandler(rewrite_handler, rewrite_func)

	s := &http.Server{
		Addr:    ":9434",
		Handler: rewrite_handler,
	}

	defer s.Close()

	go s.ListenAndServe()

	rsp, _ := http.Get("http://localhost:9434")

	defer rsp.Body.Close()

	body, _ := io.ReadAll(rsp.Body)
	log.Println(string(body))
	
	// Prints
	// <html><head><title>Test</title></head><body><p>HELLO WORLD</p></body></html>
}
```

### RewriteRequestHandler

Middleware handler to modify the `http.Request` instance passed to another handler before invoking it.

```
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func baseRequestHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		msg := req.Header.Get("X-Message")

		// See the way there is no input validation on msg? Don't do this in production...
		
		body := fmt.Sprintf("<html><head><title>Test</title></head><body><p>%s</p></body></html>", msg)
		rsp.Header().Set("Content-type", "text/html")
		rsp.Write([]byte(body))
	}

	h := http.HandlerFunc(fn)
	return h
}

func TestRewriteRequestHandler(t *testing.T) {

	rewrite_func := func(req *http.Request) (*http.Request, error) {
		req.Header.Set("X-Message", "hello world")
		return req, nil
	}

	request_handler := baseRequestHandler()

	request_handler = RewriteRequestHandler(request_handler, rewrite_func)

	s := &http.Server{
		Addr:    ":9664",
		Handler: request_handler,
	}

	defer s.Close()

	go s.ListenAndServe()

	rsp, _ := http.Get("http://localhost:9664")

	defer rsp.Body.Close()

	body, _ := io.ReadAll(rsp.Body)
	log.Println(string(body))

	// Prints:
	// <html><head><title>Test</title></head><body><p>hello world</p></body></html>
}
```
	