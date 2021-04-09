package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"syscall/js"
)

var append_func js.Func

func main() {

	append_func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		if len(args) != 1 {
			log.Println("Invalid arguments")
			return nil
		}

		b64_img := args[0].String()
		enc_props := args[1].String()

		// decode the images

		r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64_img))

		im, _, err := image.Decode(r)

		if err != nil {
			return nil
		}

		// decode the properties to append

		var props map[string]interface{}

		err := json.Unmarshal([]byte(enc_props), &props)

		if err != nil {
			return nil
		}

		// https://pkg.go.dev/github.com/dsoprea/go-png-image-structure/v2?utm_source=godoc#example-ChunkSlice.SetExif

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		opts := jpeg.Options{Quality: 100}
		err = jpeg.Encode(wr, im, &opts)

		if err != nil {
			return nil
		}

		// err = png.Encode(wr, im)

		return base64.StdEncoding.EncodeToString(buf.Bytes())
	})

	defer parse_func.Release()

	js.Global().Set("append_exif", append_func)

	c := make(chan struct{}, 0)

	log.Println("WASM EXIF appender initialized")
	<-c
}
