package bootstrap

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/aaronland/go-http-bootstrap/static"
	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-rollup"
)

// BootstrapOptions provides a list of JavaScript and CSS link to include with HTML output.
type BootstrapOptions struct {
	// A list of relative Bootstrap Javascript URLs to append as resources in HTML output.
	JS []string
	// A list of relative Bootstrap CSS URLs to append as resources in HTML output.
	CSS []string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
	RollupAssets          bool
	Prefix                string
}

// Return a *BootstrapOptions struct with default paths and URIs.
func DefaultBootstrapOptions() *BootstrapOptions {

	opts := &BootstrapOptions{
		CSS: []string{"/css/bootstrap.min.css"},
		JS:  make([]string, 0),
	}

	return opts
}

func (opts *BootstrapOptions) EnableJavascript() {
	opts.JS = append(opts.JS, "/javascript/bootstrap.bundle.min.js")
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *BootstrapOptions) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS

	if opts.RollupAssets {

		if len(opts.CSS) > 1 {

			static_opts.CSS = []string{
				"/css/bootstrap.rollup.css",
			}
		}

		if len(opts.JS) > 1 {

			static_opts.JS = []string{
				"/javascript/bootstrap.rollup.js",
			}
		}

	}

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, opts.Prefix)
}

// Append all the files in the net/http FS instance containing the embedded Bootstrap assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux, opts *BootstrapOptions) error {

	if !opts.RollupAssets {
		return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, opts.Prefix)
	}

	js_paths := make([]string, len(opts.JS))
	css_paths := make([]string, len(opts.CSS))

	for idx, path := range opts.JS {
		path = strings.TrimLeft(path, "/")
		js_paths[idx] = path
	}

	for idx, path := range opts.CSS {
		path = strings.TrimLeft(path, "/")
		css_paths[idx] = path
	}

	switch len(js_paths) {
	case 0:
		// pass
	case 1:
		err := serveSubDir(mux, opts, "javascript")

		if err != nil {
			return fmt.Errorf("Failed to append static asset handler for javascript FS, %w", err)
		}

	default:

		rollup_js_paths := map[string][]string{
			"bootstrap.rollup.js": js_paths,
		}

		rollup_js_opts := &rollup.RollupJSHandlerOptions{
			FS:    static.FS,
			Paths: rollup_js_paths,
		}

		rollup_js_handler, err := rollup.RollupJSHandler(rollup_js_opts)

		if err != nil {
			return fmt.Errorf("Failed to create rollup JS handler, %w", err)
		}

		rollup_js_uri := "/javascript/bootstrap.rollup.js"

		if opts.Prefix != "" {

			u, err := url.JoinPath(opts.Prefix, rollup_js_uri)

			if err != nil {
				return fmt.Errorf("Failed to append prefix to %s, %w", rollup_js_uri, err)
			}

			rollup_js_uri = u
		}

		mux.Handle(rollup_js_uri, rollup_js_handler)
	}

	// CSS

	switch len(css_paths) {
	case 0:
		// pass
	case 1:

		err := serveSubDir(mux, opts, "css")

		if err != nil {
			return fmt.Errorf("Failed to append static asset handler for css FS, %w", err)
		}

	default:

		rollup_css_paths := map[string][]string{
			"bootstrap.rollup.css": css_paths,
		}

		rollup_css_opts := &rollup.RollupCSSHandlerOptions{
			FS:    static.FS,
			Paths: rollup_css_paths,
		}

		rollup_css_handler, err := rollup.RollupCSSHandler(rollup_css_opts)

		if err != nil {
			return fmt.Errorf("Failed to create rollup CSS handler, %w", err)
		}

		rollup_css_uri := "/css/bootstrap.rollup.css"

		if opts.Prefix != "" {

			u, err := url.JoinPath(opts.Prefix, rollup_css_uri)

			if err != nil {
				return fmt.Errorf("Failed to append prefix to %s, %w", rollup_css_uri, err)
			}

			rollup_css_uri = u
		}

		mux.Handle(rollup_css_uri, rollup_css_handler)
	}

	// END OF this should eventually be made a generic function in go-http-rollup

	return nil
}

func serveSubDir(mux *http.ServeMux, opts *BootstrapOptions, dirname string) error {

	sub_fs, err := fs.Sub(static.FS, dirname)

	if err != nil {
		return fmt.Errorf("Failed to load %s FS, %w", dirname, err)
	}

	sub_prefix := dirname

	if opts.Prefix != "" {

		prefix, err := url.JoinPath(opts.Prefix, sub_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append prefix to %s, %w", sub_prefix, err)
		}

		sub_prefix = prefix
	}

	err = aa_static.AppendStaticAssetHandlersWithPrefix(mux, sub_fs, sub_prefix)

	if err != nil {
		return fmt.Errorf("Failed to append static asset handler for %s FS, %w", dirname, err)
	}

	return nil
}
