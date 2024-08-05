// `go-http-bootstrap` is an HTTP middleware package for including Bootstrap.js assets in web applications. It exports two principal methods:
//
// * `bootstrap.AppendAssetHandlers(*http.ServeMux)` which is used to append HTTP handlers to a `http.ServeMux` instance for serving Bootstrap CSS and JavaScript files, and related assets.
// * `bootstrap.AppendResourcesHandler(http.Handler, *BootstrapOptions)` which is used to rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap JavaScript files and related assets.
//
// Example
//
//	package main
//
//	import (
//		"github.com/aaronland/go-http-bootstrap"
//		"log"
//		"net/http"
//	)
//
//	func Handler() http.Handler {
//
//		index := `
//	<!doctype html>
//	<html lang="en-us">
//	  <head>
//	    <meta charset="utf-8">
//	    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
//	    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
//	    <title>Bootstrap</title>
//	  </head>
//
//	  <body>
//	   <div class="card">
//	   	<h1 class="card-header">Card header</h1>
//		<div class="card-body">Card body</div>
//		<div class="card-footer">Card footer</div>
//	   </div>
//	  </body>
//	</html>`
//
//		fn := func(rsp http.ResponseWriter, req *http.Request) {
//
//			rsp.Write([]byte(index))
//		}
//
//		return http.HandlerFunc(fn)
//	}
//
//	func main() {
//
//		mux := http.NewServeMux()
//
//		idx_handler := Handler()
//
//		bootstrap_opts := bootstrap.DefaultBootstrapOptions()
//		idx_handler = bootstrap.AppendResourcesHandler(idx_handler, bootstrap_opts)
//
//		mux.Handle("/", idx_handler)
//
//		bootstrap.AppendAssetHandlers(mux)
//
//		endpoint := "localhost:8080"
//		log.Printf("Listening for requests on %s\n", endpoint)
//
//		http.ListenAndServe(endpoint, mux)
//	}
//
// All of the Bootstrap files in the [static/css](static/css) and [static/javascript](static/javascript) are registered with your `http.ServeMux` instance when you call `bootstrap.AppendAssetHandlers` but by default only the `css/bootstrap.min.css` is included in the list of CSS and Javascript resources to append to HTML content when you call the `bootstrap.DefaultBootstrapOptions()` method. If there are other Bootstrap-related files you need to access in your application you will need to add them to the `BootstrapOptions.CSS` and `Bootstrap.JS` properties manually.
package bootstrap
