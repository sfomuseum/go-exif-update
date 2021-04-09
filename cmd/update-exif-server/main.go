package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-exif-wasm/www"
	"log"
	"net/http"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

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

	http_fs := http.FS(www.FS)
	fs_handler := http.FileServer(http_fs)

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()
	fs_handler = bootstrap.AppendResourcesHandler(fs_handler, bootstrap_opts)

	mux.Handle("/", fs_handler)

	log.Printf("Listening on %s", s.Address())
	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
