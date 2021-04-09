// package update provides methods for updating EXIF data in JPEG files.
package update

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"io"
	"log"
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

	/*
		ifdPath := exifcommon.Ifd1StandardIfdIdentity

		ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath.Name())

		if err != nil {
			return err
		}
	*/

	// https://github.com/dsoprea/go-exif/blob/de2141190595193aa097a2bf3205ba0cf76dc14b/tags_data.go

	paths := []*exifcommon.IfdIdentity{
		exifcommon.IfdStandardIfdIdentity,
		exifcommon.IfdExifStandardIfdIdentity,
		exifcommon.IfdExifIopStandardIfdIdentity,
		exifcommon.IfdGpsInfoStandardIfdIdentity,
		exifcommon.Ifd1StandardIfdIdentity,
	}

	ti := exif.NewTagIndex()

	for k, v := range exif_props {

		for _, p := range paths {

			_, err := ti.GetWithName(p, k)

			if err != nil {
				continue
			}

			err = setExifTag(rootIb, p.UnindexedString(), k, v.(string))

			if err != nil {
				return nil
			}

			log.Println("SET", p, k, v, err)
		}
	}

	/*
		for k, v := range exif_props {


			err = ifdIb.SetStandardWithName(k, v) //"CameraOwnerName", "SFO Museum")

			if err != nil {
				return fmt.Errorf("Failed to set property '%s', %v", k, err)
			}
		}
	*/

	// Update the exif segment.

	err = sl.SetExif(rootIb)

	if err != nil {
		return err
	}

	return sl.Write(wr)
}

func setExifTag(rootIB *exif.IfdBuilder, ifdPath, tagName, tagValue string) error {

	fmt.Printf("setTag(): ifdPath: %v, tagName: %v, tagValue: %v",
		ifdPath, tagName, tagValue)

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIB, ifdPath)

	if err != nil {
		return fmt.Errorf("Failed to get or create IB for %s: %v", ifdPath, err)
	}

	err = ifdIb.SetStandardWithName(tagName, tagValue)

	if err != nil {
		return fmt.Errorf("failed to set %s tag: %v", tagName, err)
	}

	return nil
}
