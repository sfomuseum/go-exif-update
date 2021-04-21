// package update provides methods for updating EXIF data in JPEG files.
package update

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/sfomuseum/go-exif-update/tags"
	"io"
	"log"
)

var ti *exif.TagIndex

var tag_paths []*exifcommon.IfdIdentity

func init() {

	ti = exif.NewTagIndex()

	// https://github.com/dsoprea/go-exif/blob/de2141190595193aa097a2bf3205ba0cf76dc14b/tags_data.go

	tag_paths = []*exifcommon.IfdIdentity{
		exifcommon.IfdStandardIfdIdentity,
		exifcommon.IfdExifStandardIfdIdentity,
		exifcommon.IfdExifIopStandardIfdIdentity,
		exifcommon.IfdGpsInfoStandardIfdIdentity,
		exifcommon.Ifd1StandardIfdIdentity,
	}

}

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

	for k, v := range exif_props {

		ok, err := tags.IsSupported(k)

		if err != nil {
			return err
		}

		if !ok {
			return fmt.Errorf("Tag '%s' is not supported at this time", k)
		}

		id, _, err := GetIndexedTagFromName(k)

		if err != nil {
			return err
		}

		log.Println("SET", k, v)

		err = setExifTag(rootIb, id.UnindexedString(), k, v)

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
}

// Cribbed from https://github.com/dsoprea/go-exif/issues/11

func setExifTag(rootIB *exif.IfdBuilder, ifdPath string, tagName string, tagValue interface{}) error {

	// log.Printf("setTag(): ifdPath: %v, tagName: %v, tagValue: %v", ifdPath, tagName, tagValue)

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

// Return the *exifcommon.IfdIdentity and *exif.IndexedTag instances associated
// with a given EXIF string tag name.
func GetIndexedTagFromName(k string) (*exifcommon.IfdIdentity, *exif.IndexedTag, error) {

	for _, id := range tag_paths {

		t, err := ti.GetWithName(id, k)

		if err != nil {
			continue
		}

		return id, t, nil
	}

	return nil, nil, fmt.Errorf("Unrecognized tag, %s", k)
}
