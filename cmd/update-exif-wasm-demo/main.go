package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-exif-update/www"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-http-wasm"		
)

func main() {

	fs := flagset.NewFlagSet("server")

	server_uri := fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	bootstrap_prefix := fs.String("bootstrap-prefix", "", "A relative path to append to all Bootstrap-related paths the server will listen for requests on.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "HTTP server for demonstrating the use of the update_exif WebAssembly binary.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "EXIF", true)

	if err != nil {
		log.Fatalf("Failed to set flags from environment variables, %v", err)
	}

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	mux := http.NewServeMux()

	err = bootstrap.AppendAssetHandlers(mux)

	if err != nil {
		log.Fatalf("Failed to append Bootstrap asset handlers, %v", err)
	}

	err = wasm.AppendAssetHandlers(mux)

	if err != nil {
		log.Fatalf("Failed to append WASM asset handlers, %v", err)
	}
	
	http_fs := http.FS(www.FS)
	fs_handler := http.FileServer(http_fs)

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()
	fs_handler = bootstrap.AppendResourcesHandlerWithPrefix(fs_handler, bootstrap_opts, *bootstrap_prefix)

	wasm_opts := wasm.DefaultWASMOptions()
	fs_handler = wasm.AppendResourcesHandler(fs_handler, wasm_opts)
	
	mux.Handle("/", fs_handler)

	log.Printf("Listening on %s", s.Address())
	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
