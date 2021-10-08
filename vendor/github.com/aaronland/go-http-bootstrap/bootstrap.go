package bootstrap

import (
	"fmt"
	"github.com/aaronland/go-http-rewrite"	
	"github.com/aaronland/go-http-bootstrap/static"
	"io/fs"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

// BootstrapOptions provides a list of JavaScript and CSS link to include with HTML output.
type BootstrapOptions struct {
	// A list of relative Bootstrap Javascript URLs to append as resources in HTML output.
	JS  []string
	// A list of relative Bootstrap CSS URLs to append as resources in HTML output.	
	CSS []string
}

// Return a *BootstrapOptions struct with default paths and URIs.
func DefaultBootstrapOptions() *BootstrapOptions {

	opts := &BootstrapOptions{
		CSS: []string{"/css/bootstrap.min.css"},
		JS:  make([]string, 0),
	}

	return opts
}

// AppendResourcesHandler will rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap JavaScript files and related assets.
func AppendResourcesHandler(next http.Handler, opts *BootstrapOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Bootstrap JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandlerWithPrefix(next http.Handler, opts *BootstrapOptions, prefix string) http.Handler {

	// We're doing this the long way because otherwise there is a
	// risk of infinite-prefixing because of copy by reference issues
	// (20210322/straup)

	js := make([]string, len(opts.JS))
	css := make([]string, len(opts.CSS))

	for idx, path := range opts.JS {

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		js[idx] = path
	}

	for idx, path := range opts.CSS {

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		css[idx] = path
	}

	rewrite_opts := &rewrite.AppendResourcesOptions{
		JavaScript:  js,
		Stylesheets: css,
	}

	return rewrite.AppendResourcesHandler(next, rewrite_opts)
}

// AssetsHandler returns a net/http FS instance containing the embedded Bootstrap assets that are included with this package.
func AssetsHandler() (http.Handler, error) {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

// AssetsHandler returns a net/http FS instance containing the embedded Bootstrap assets that are included with this package ensuring that all URLs are stripped of prefix.
func AssetsHandlerWithPrefix(prefix string) (http.Handler, error) {

	fs_handler, err := AssetsHandler()

	if err != nil {
		return nil, err
	}

	fs_handler = http.StripPrefix(prefix, fs_handler)
	return fs_handler, nil
}

// Append all the files in the net/http FS instance containing the embedded Bootstrap assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

// Append all the files in the net/http FS instance containing the embedded Bootstrap assets to an *http.ServeMux instance ensuring that all URLs are prepended with prefix.
func AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	asset_handler, err := AssetsHandlerWithPrefix(prefix)

	if err != nil {
		return nil
	}

	walk_func := func(path string, info fs.DirEntry, err error) error {

		if path == "." {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		mux.Handle(path, asset_handler)
		return nil
	}

	return fs.WalkDir(static.FS, ".", walk_func)
}

func appendPrefix(prefix string, path string) string {

	prefix = strings.TrimRight(prefix, "/")

	if prefix != "" {
		path = strings.TrimLeft(path, "/")
		path = filepath.Join(prefix, path)
	}

	return path
}
