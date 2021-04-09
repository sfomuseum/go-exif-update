package main

import (
	"flag"
	"github.com/sfomuseum/go-exif-update"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"os"
)

func main() {

	var properties multi.KeyValueString
	flag.Var(&properties, "property", "A {TAG}={VALUE} EXIF property string. {TAG} must be a recognized EXIF tag.")

	flag.Parse()

	paths := flag.Args()

	exif_props := make(map[string]interface{})

	for _, p := range properties {
		k := p.Key()
		v := p.Value().(string)

		exif_props[k] = v
	}

	for _, path := range paths {

		fh, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", path, err)
		}

		defer fh.Close()

		err = update.UpdateExif(fh, os.Stdout, exif_props)

		if err != nil {
			log.Fatalf("Failed to update '%s', %v", path, err)
		}
	}

}
