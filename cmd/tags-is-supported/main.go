package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sfomuseum/go-exif-update/tags"	
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for indicating whether a named EXIF tag is supported by the sfomuseum/go-exif-update package.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s tag(N) tag(N) tag(N)\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	for _, t := range flag.Args() {

		ok, err := tags.IsSupported(t)

		if err != nil {
			log.Fatalf("Failed to determine whether tag '%s' is supported, %v", t, err)
		}

		fmt.Printf("%s %t\n", t, ok)
	}
}
