package update

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"io"
)

func UpdateExifB64(r io.Reader, wr io.Writer, exif_data map[string]interface{}) error {

	img_fh := base64.NewDecoder(base64.StdEncoding, r)

	var buf bytes.Buffer
	img_wr := bufio.NewWriter(&buf)

	err := UpdateExif(img_fh, img_wr, exif_data)

	if err != nil {
		return err
	}

	img_wr.Flush()

	enc := base64.NewEncoder(base64.StdEncoding, wr)

	_, err = enc.Write(buf.Bytes())

	if err != nil {
		return err
	}

	return enc.Close()
}

func UpdateExif(r io.Reader, wr io.Writer, exif_data map[string]interface{}) error {

	img_data, err := io.ReadAll(r)

	// https://pkg.go.dev/github.com/dsoprea/go-jpeg-image-structure/v2?utm_source=godoc#example-SegmentList.SetExif

	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(img_data)

	if err != nil {
		return err
	}

	sl := intfc.(*jpegstructure.SegmentList)

	rootIb, err := sl.ConstructExifBuilder()

	if err != nil {
		return err
	}

	ifdPath := "IFD/Exif"

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)

	if err != nil {
		return err
	}

	for k, v := range exif_data {

		err = ifdIb.SetStandardWithName(k, v) //"CameraOwnerName", "SFO Museum")

		if err != nil {
			return err
		}
	}

	// Update the exif segment.

	err = sl.SetExif(rootIb)

	if err != nil {
		return err
	}

	return sl.Write(wr)

	/*

		d := buf.Bytes()

		intfc, err = jmp.ParseBytes(d)

		if err != nil {
			return nil, err
		}

		sl = intfc.(*jpegstructure.SegmentList)

		_, _, exifTags, err := sl.DumpExif()

		if err != nil {
			return nil, err
		}

		for _, et := range exifTags {
			if et.IfdPath == "IFD/Exif" && et.TagName == "CameraOwnerName" {
				fmt.Printf("Value: [%s]\n", et.FormattedFirst)
				break
			}
		}

		return nil, nil
	*/
}
