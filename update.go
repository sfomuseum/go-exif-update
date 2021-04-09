// package update provides methods for updating EXIF data in JPEG files.
package update

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"io"
)

// UpdateExif updates the EXIF data encoded in r and writes that data to wr.
// This is really nothing more than a thin wrapper around the example code in
// dsoprea's go-jpeg-image-structure package.
func UpdateExif(r io.Reader, wr io.Writer, exif_props map[string]interface{}) error {

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

	for k, v := range exif_props {

		err = ifdIb.SetStandardWithName(k, v) //"CameraOwnerName", "SFO Museum")

		if err != nil {
			return fmt.Errorf("Failed to set property '%s', %v", k, err)
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
