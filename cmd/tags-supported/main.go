package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/sfomuseum/go-exif-update/tags"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool that prints a list of EXIF tag names, sorted alphabetically, that are supported by the sfomuseum/go-exif-update package.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	tags, err := tags.SupportedTags()

	if err != nil {
		log.Fatalf("Failed to derive list of supported tags, %v", err)
	}

	sort.Strings(tags)

	for _, t := range tags {
		fmt.Println(t)
	}
}
