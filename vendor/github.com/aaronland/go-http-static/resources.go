package static

import (
	"net/http"

	"github.com/aaronland/go-http-rewrite"
)

// ResourcesOptions provides a list of JavaScript and CSS link to include with HTML output.
type ResourcesOptions struct {
	JS                    []string
	CSS                   []string
	DataAttributes        map[string]string
	AppendJavaScriptAtEOF bool
}

// Return a *ResourcesOptions struct with default paths and URIs.
func DefaultResourcesOptions() *ResourcesOptions {

	opts := &ResourcesOptions{
		CSS:            []string{},
		JS:             []string{},
		DataAttributes: make(map[string]string),
	}

	return opts
}

func AppendResourcesHandler(next http.Handler, opts *ResourcesOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

func AppendResourcesHandlerWithPrefix(next http.Handler, opts *ResourcesOptions, prefix string) http.Handler {

	js := make([]string, len(opts.JS))
	css := make([]string, len(opts.CSS))

	for i, path := range opts.JS {
		js[i] = appendPrefix(prefix, path)
	}

	for i, path := range opts.CSS {
		css[i] = appendPrefix(prefix, path)
	}

	ext_opts := &rewrite.AppendResourcesOptions{
		JavaScript:            js,
		Stylesheets:           css,
		DataAttributes:        opts.DataAttributes,
		AppendJavaScriptAtEOF: opts.AppendJavaScriptAtEOF,
	}

	return rewrite.AppendResourcesHandler(next, ext_opts)
}
