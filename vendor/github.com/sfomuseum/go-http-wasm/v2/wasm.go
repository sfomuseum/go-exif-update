package wasm

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-rollup"
	"github.com/sfomuseum/go-http-wasm/v2/static"
)

// WASMOptions provides a list of JavaScript and CSS link to include with HTML output.
type WASMOptions struct {
	JS []string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
	RollupAssets          bool
	Prefix                string
}

// Return a *WASMOptions struct with default paths and URIs.
func DefaultWASMOptions() *WASMOptions {

	opts := &WASMOptions{
		JS: []string{
			"/javascript/wasm_exec.js",
			"/javascript/sfomuseum.wasm.js",
		},
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load WASM JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *WASMOptions) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.JS = opts.JS
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	if opts.RollupAssets {

		static_opts.JS = []string{
			"/javascript/sfomuseum.wasm.rollup.js",
		}
	}

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, opts.Prefix)
}

// Append all the files in the net/http FS instance containing the embedded WASM assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux, opts *WASMOptions) error {

	if !opts.RollupAssets {
		return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, opts.Prefix)
	}

	js_paths := make([]string, len(opts.JS))

	for idx, path := range opts.JS {
		path = strings.TrimLeft(path, "/")
		js_paths[idx] = path
	}

	rollup_js_paths := map[string][]string{
		"sfomuseum.wasm.rollup.js": js_paths,
	}

	rollup_js_opts := &rollup.RollupJSHandlerOptions{
		FS:    static.FS,
		Paths: rollup_js_paths,
	}

	rollup_js_handler, err := rollup.RollupJSHandler(rollup_js_opts)

	if err != nil {
		return fmt.Errorf("Failed to create rollup JS handler, %w", err)
	}

	rollup_js_uri := "/javascript/sfomuseum.wasm.rollup.js"

	if opts.Prefix != "" {

		u, err := url.JoinPath(opts.Prefix, rollup_js_uri)

		if err != nil {
			return fmt.Errorf("Failed to append prefix to %s, %w", rollup_js_uri, err)
		}

		rollup_js_uri = u
	}

	mux.Handle(rollup_js_uri, rollup_js_handler)
	return nil
}
