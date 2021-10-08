# go-http-bootstrap

`go-http-bootstrap` is an HTTP middleware package for including Bootstrap (v5.0.0) assets in web applications.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-http-bootstrap.svg)](https://pkg.go.dev/github.com/aaronland/go-http-bootstrap)

`go-http-bootstrap` is an HTTP middleware package for including Bootstrap.js assets in web applications. It exports two principal methods:

* `bootstrap.AppendAssetHandlers(*http.ServeMux)` which is used to append HTTP handlers to a `http.ServeMux` instance for serving Bootstrap CSS and JavaScript files, and related assets.
* `bootstrap.AppendResourcesHandler(http.Handler, *BootstrapOptions)` which is used to rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap

This package doesn't specify any code or methods for how Bootstrap.js is used. It only provides method for making Bootstraps available to existing applications.

## Example

```
package main

import (
	"github.com/aaronland/go-http-bootstrap"
	"log"
	"net/http"
)

func Handler() http.Handler {

	index := `
<!doctype html>
<html lang="en-us">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>Bootstrap</title>
  </head>

  <body>
   <div class="card">
   	<h1 class="card-header">Card header</h1>
	<div class="card-body">Card body</div>
	<div class="card-footer">Card footer</div>
   </div>
  </body>
</html>`

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Write([]byte(index))
	}

	return http.HandlerFunc(fn)
}

func main() {

	mux := http.NewServeMux()
	
	idx_handler := Handler()

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()
	idx_handler = bootstrap.AppendResourcesHandler(idx_handler, bootstrap_opts)

	mux.Handle("/", idx_handler)

	bootstrap.AppendAssetHandlers(mux)

	endpoint := "localhost:8080"
	log.Printf("Listening for requests on %s\n", endpoint)

	http.ListenAndServe(endpoint, mux)
}
```

_Error handling omitted for brevity._

You can see an example of this application by running the [cmd/example](cmd/example/main.go) application. You can do so by invoking the `example` Makefile target. For example:

```
$> make example
go run -mod vendor cmd/example/main.go 
2021/05/05 13:54:07 Listening for requests on localhost:8080
```

The when you open the URL `http://localhost:8080` in a web browser you should see the following:

![](docs/images/go-http-bootstrap-example.png)

### Notes

All of the Bootstrap files in the [static/css](static/css) and [static/javascript](static/javascript) are registered with your `http.ServeMux` instance when you call `bootstrap.AppendAssetHandlers` but by default only the `css/bootstrap.min.css` is included in the list of CSS and Javascript resources to append to HTML content when you call the `bootstrap.DefaultBootstrapOptions()` method. If there are other Bootstrap-related files you need to access in your application you will need to add them to the `BootstrapOptions.CSS` and `Bootstrap.JS` properties manually.

## See also

* https://getbootstrap.com/
