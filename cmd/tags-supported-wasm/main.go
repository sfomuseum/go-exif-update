package main

import (
	"encoding/json"
	"github.com/sfomuseum/go-exif-update/tags"
	"log"
	"syscall/js"
)

var supported_func js.Func

func main() {

	supported_func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		tags_supported, err := tags.SupportedTags()

		if err != nil {
			log.Printf("Failed to derive supported tags, %v", err)
			return nil
		}

		enc_supported, err := json.Marshal(tags_supported)

		if err != nil {
			log.Printf("Failed to encode supported tags, %v", err)
			return nil
		}

		return string(enc_supported)
	})

	defer supported_func.Release()

	js.Global().Set("supported_tags", supported_func)

	c := make(chan struct{}, 0)

	log.Println("WASM EXIF supported tags initialized")
	<-c
}
