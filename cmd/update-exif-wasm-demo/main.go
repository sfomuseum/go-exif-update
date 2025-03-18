package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-exif-update/www"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	fs := flagset.NewFlagSet("server")

	server_uri := fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

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

	http_fs := http.FS(www.FS)
	fs_handler := http.FileServer(http_fs)

	mux.Handle("/", fs_handler)

	log.Printf("Listening on %s", s.Address())
	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
