package main

import (
	"flag"
	"github.com/sfomuseum/go-exif-wasm"
	"log"
	"os"
)

func main() {

	flag.Parse()

	paths := flag.Args()

	first := paths[0]
	fh, err := os.Open(first)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", first, err)
	}

	defer fh.Close()

	exif_data := map[string]interface{}{
		"CameraOwnerName": "SFO Museum, yo",
	}

	err = update.UpdateExif(fh, os.Stdout, exif_data)

	if err != nil {
		log.Fatal(err)
	}
}
