package bootstrap

import (
	"fmt"
	"github.com/aaronland/go-http-bootstrap/resources"
	"github.com/aaronland/go-http-bootstrap/static"
	"io/fs"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

type BootstrapOptions struct {
	JS  []string
	CSS []string
}

func DefaultBootstrapOptions() *BootstrapOptions {

	opts := &BootstrapOptions{
		CSS: []string{"/css/bootstrap.min.css"},
		JS:  make([]string, 0),
	}

	return opts
}

func AppendResourcesHandler(next http.Handler, opts *BootstrapOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

func AppendResourcesHandlerWithPrefix(next http.Handler, opts *BootstrapOptions, prefix string) http.Handler {

	// We're doing this the long way because otherwise there is a
	// risk of infinite-prefixing because of copy by reference issues
	// (20210322/straup)

	ext_opts := &resources.AppendResourcesOptions{
		JS:  make([]string, len(opts.JS)),
		CSS: make([]string, len(opts.CSS)),
	}

	for idx, path := range opts.JS {

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		ext_opts.JS[idx] = path
	}

	for idx, path := range opts.CSS {

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		ext_opts.CSS[idx] = path
	}

	return resources.AppendResourcesHandler(next, ext_opts)
}

func AssetsHandler() (http.Handler, error) {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

func AssetsHandlerWithPrefix(prefix string) (http.Handler, error) {

	fs_handler, err := AssetsHandler()

	if err != nil {
		return nil, err
	}

	fs_handler = http.StripPrefix(prefix, fs_handler)
	return fs_handler, nil
}

func AppendAssetHandlers(mux *http.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

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
