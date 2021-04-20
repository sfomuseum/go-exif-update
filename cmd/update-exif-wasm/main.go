package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-exif-update"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"regexp"
	"strings"
	"syscall/js"
)

// var update_func js.Func
var b64_pat *regexp.Regexp

func init() {
	b64_pat = regexp.MustCompile(`^data:image/(\w+);base64,(.*)$`)
}

func UpdateFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		b64_data := args[0].String()
		enc_props := args[1].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			log.Println("WOO", resolve, reject)

			go func() {

				log.Println("GOING...")

				if !b64_pat.MatchString(b64_data) {
					log.Println("SAD", 1)
					reject.Invoke("Base64 image data failed match")
					return
				}

				log.Println("STEP", 1)
				m := b64_pat.FindStringSubmatch(b64_data)

				format := m[1]
				b64_img := m[2]

				// decode the EXIF properties

				log.Println("STEP", 2)

				var exif_data map[string]interface{}

				err := json.Unmarshal([]byte(enc_props), &exif_data)

				if err != nil {
					log.Println("SAD", 2)
					reject.Invoke(fmt.Sprintf("Failed to unmarshal properties, %v", err))
					return
				}

				log.Println("STEP", 3)
				// decode the images

				b64_fh := strings.NewReader(b64_img)
				var img_fh io.Reader

				log.Println("STEP", 4)
				if format == "jpeg" {
					img_fh = base64.NewDecoder(base64.StdEncoding, b64_fh)
				} else {

					tmp_fh := base64.NewDecoder(base64.StdEncoding, b64_fh)

					im, _, err := image.Decode(tmp_fh)

					if err != nil {
						log.Println("SAD", 3)
						reject.Invoke(fmt.Sprintf("Failed to decode image data, %v", err))
						return
					}

					jpg_r, jpg_wr := io.Pipe()

					opts := jpeg.Options{Quality: 100}
					err = jpeg.Encode(jpg_wr, im, &opts)

					if err != nil {
						log.Println("SAD", 4)
						reject.Invoke(fmt.Sprintf("Failed to decode image data as JPEG, %v", err))
						return
					}

					img_fh = jpg_r
				}

				log.Println("STEP", 5)
				var buf bytes.Buffer
				img_wr := bufio.NewWriter(&buf)

				err = update.UpdateExif(img_fh, img_wr, exif_data)

				if err != nil {
					log.Println("SAD", 5)
					reject.Invoke(fmt.Printf("Failed update EXIF properties, %v", err))
					return
				}

				log.Println("STEP", 6)
				img_wr.Flush()

				b64_img = base64.StdEncoding.EncodeToString(buf.Bytes())
				data_uri := fmt.Sprintf("data:image/jpeg;base64,%s", b64_img)

				log.Println("RESOLVE OKAY", data_uri)
				resolve.Invoke(data_uri)
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func main() {

	update_func := UpdateFunc()
	defer update_func.Release()

	js.Global().Set("update_exif", update_func)

	c := make(chan struct{}, 0)

	log.Println("WASM EXIF updater initialized")
	<-c
}
