package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-exif-update"
	"image"
	"image/jpeg"
	_ "image/png"
	_ "image/gif"	
	"io"
	"log"
	"regexp"
	"strings"
	"syscall/js"
)

var update_func js.Func
var b64_pat *regexp.Regexp

func main() {

	b64_pat, err := regexp.Compile(`^data:image/(\w+);base64,(.*)$`)

	if err != nil {
		log.Fatalf("Failed to compile B64 pattern, %v", err)
	}

	update_func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		if len(args) != 2 {
			log.Println("Invalid arguments")
			return nil
		}

		b64_data := args[0].String()
		enc_props := args[1].String()

		//

		if !b64_pat.MatchString(b64_data) {
			return nil
		}

		m := b64_pat.FindStringSubmatch(b64_data)

		format := m[1]
		b64_img := m[2]

		// decode the EXIF properties

		var exif_data map[string]interface{}

		err := json.Unmarshal([]byte(enc_props), &exif_data)

		if err != nil {
			return nil
		}

		// decode the images

		b64_fh := strings.NewReader(b64_img)
		var img_fh io.Reader

		if format == "jpeg" {
			img_fh = base64.NewDecoder(base64.StdEncoding, b64_fh)
		} else {

			tmp_fh := base64.NewDecoder(base64.StdEncoding, b64_fh)

			im, _, err := image.Decode(tmp_fh)

			if err != nil {
				return nil
			}

			jpg_r, jpg_wr := io.Pipe()

			opts := jpeg.Options{Quality: 100}
			err = jpeg.Encode(jpg_wr, im, &opts)

			if err != nil {
				return err
			}

			img_fh = jpg_r
		}

		var buf bytes.Buffer
		img_wr := bufio.NewWriter(&buf)

		err = update.UpdateExif(img_fh, img_wr, exif_data)

		if err != nil {
			return err
		}

		img_wr.Flush()

		b64_img = base64.StdEncoding.EncodeToString(buf.Bytes())
		return fmt.Sprintf("data:image/jpeg;base64,%s", b64_img)
	})

	defer update_func.Release()

	js.Global().Set("update_exif", update_func)

	c := make(chan struct{}, 0)

	log.Println("WASM EXIF updater initialized")
	<-c
}
