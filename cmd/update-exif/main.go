package main

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-exif-update"
	"github.com/sfomuseum/go-exif-update/tags"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"os"
	"strings"
)

func main() {

	var properties multi.KeyValueString
	flag.Var(&properties, "property", "One or more {TAG}={VALUE} EXIF properties. {TAG} must be a recognized EXIF tag.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for updating the EXIF properties in one or more JPEG images. Images are not updated in place but written to STDOUT.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] image(N) image(N) image(N)\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	paths := flag.Args()

	exif_props := make(map[string]interface{})

	for _, p := range properties {

		k := p.Key()
		v := p.Value().(string)

		if !strings.HasPrefix(k, "X-") {

			ok, err := tags.IsSupported(k)

			if err != nil {
				log.Fatalf("Failed to determine whether tag '%s' is supported, %v", k, err)
			}

			if !ok {
				log.Fatalf("Tag '%s' is not supported by this tool, at this time", k)
			}
		}

		exif_props[k] = v
	}

	for _, path := range paths {

		fh, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", path, err)
		}

		defer fh.Close()

		err = update.PrepareAndUpdateExif(fh, os.Stdout, exif_props)

		if err != nil {
			log.Fatalf("Failed to update '%s', %v", path, err)
		}
	}

}
