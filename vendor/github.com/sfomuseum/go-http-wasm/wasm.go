package wasm

import (
	gohttp "net/http"

	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-wasm/static"
)

// WASMOptions provides a list of JavaScript and CSS link to include with HTML output.
type WASMOptions struct {
	JS  []string
}

// Return a *WASMOptions struct with default paths and URIs.
func DefaultWASMOptions() *WASMOptions {

	opts := &WASMOptions{
		JS: []string{
			"/javascript/wasm_exec.js",
		},
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load WASM JavaScript files and related assets.
func AppendResourcesHandler(next gohttp.Handler, opts *WASMOptions) gohttp.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load WASM JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next gohttp.Handler, opts *WASMOptions, prefix string) gohttp.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.JS = opts.JS

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}

// Append all the files in the net/http FS instance containing the embedded WASM assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *gohttp.ServeMux) error {

	return aa_static.AppendStaticAssetHandlers(mux, static.FS)
}

// Append all the files in the net/http FS instance containing the embedded WASM assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *gohttp.ServeMux, prefix string) error {

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
