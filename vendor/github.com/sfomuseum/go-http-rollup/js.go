package rollup

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"io/fs"
	"log/slog"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

type RollupJSHandlerOptions struct {
	FS fs.FS
	Paths map[string][]string
}

func RollupJSHandler(opts *RollupJSHandlerOptions) (http.Handler, error) {

	js_regexp, err := regexp.Compile("^(application|text)/(x-)?(java|ecma)script$")

	if err != nil {
		return nil, fmt.Errorf("Failed to compile JS pattern, %w", err)
	}

	m := minify.New()
	m.AddFuncRegexp(js_regexp, js.Minify)

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		path := req.URL.Path
		fname := filepath.Base(path)

		js_paths, exists := opts.Paths[fname]

		if !exists {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		rsp.Header().Set("Content-type", "text/javascript")

		for _, path := range js_paths {

			r, err := opts.FS.Open(path)

			if err != nil {
				slog.Error("Failed to open JavaScript file for reading", "path", path, "error", err)
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			defer r.Close()

			err = m.Minify("text/javascript", rsp, r)

			if err != nil {
				slog.Error("Failed to minify JavaScript file", "path", path, "error", err)				
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			rsp.Write([]byte(`;`))
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
