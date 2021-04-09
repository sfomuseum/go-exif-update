package main

import (
	"fmt"
	"github.com/sfomuseum/go-exif-update/tags"
	"log"
	"sort"
)

func main() {

	tags, err := tags.SupportedTags()

	if err != nil {
		log.Fatalf("Failed to derive list of supported tags, %v", err)
	}

	sort.Strings(tags)

	for _, t := range tags {
		fmt.Println(t)
	}
}
