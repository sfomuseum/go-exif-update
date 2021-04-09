package main

import (
	_ "bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-png-image-structure/v2"
	_ "image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
)

func append(fh io.Reader) (string, error) {

	img_fh := base64.NewDecoder(base64.StdEncoding, fh)

	img_data, err := io.ReadAll(img_fh)

	if err != nil {
		return "", err
	}

	// Create EXIF

	im, err := exifcommon.NewIfdMappingWithStandard()

	if err != nil {
		return "", err
	}

	ti := exif.NewTagIndex()

	ib := exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.TestDefaultByteOrder)

	err = ib.AddStandardWithName("ImageWidth", []uint32{11})

	if err != nil {
		return "", err
	}

	// Update PNG file

	pmp := pngstructure.NewPngMediaParser()

	intfc, err := pmp.ParseBytes(img_data)

	if err != nil {
		return "", err
	}

	cs := intfc.(*pngstructure.ChunkSlice)

	err = cs.SetExif(ib)

	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	err = cs.WriteTo(b)

	if err != nil {
		return "", err
	}

	return "", nil

	/*
		r := base64.NewDecoder(base64.StdEncoding, fh)

		im, _, err := image.Decode(r)

		if err != nil {
			return "", err
		}

		// https://pkg.go.dev/github.com/dsoprea/go-png-image-structure/v2?utm_source=godoc#example-ChunkSlice.SetExif

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		opts := jpeg.Options{Quality: 100}
		err = jpeg.Encode(wr, im, &opts)

		if err != nil {
			return "", err
		}

		// err = png.Encode(wr, im)

		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	*/
}

func main() {

	flag.Parse()

	paths := flag.Args()

	first := paths[0]
	fh, err := os.Open(first)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", first, err)
	}

	defer fh.Close()

	new, err := append(fh)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(new)
}
