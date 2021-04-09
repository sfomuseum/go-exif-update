package main

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-exif-update/tags"
	"log"
)

func main() {

	flag.Parse()

	for _, t := range flag.Args() {

		ok, err := tags.IsSupported(t)

		if err != nil {
			log.Fatalf("Failed to determine whether tag '%s' is supported, %v", t, err)
		}

		fmt.Printf("%s %t\n", t, ok)
	}
}
