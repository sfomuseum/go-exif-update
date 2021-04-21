package update

import (
	"io"
	"os"
	"testing"
)

func TestUpdateExifStrings(t *testing.T) {

	r, err := os.Open("fixtures/walrus.jpg")

	if err != nil {
		t.Fatalf("Failed to open test image, %v", err)
	}

	defer r.Close()

	wr := io.Discard

	props := map[string]interface{}{
		"Artist": "Bob",
	}

	err = UpdateExif(r, wr, props)

	if err != nil {
		t.Fatalf("Failed to update EXIF data, %v", err)
	}

	// TO DO: READ AND VALIDATE TAGS
}

func TestUpdateExifRationals(t *testing.T) {

	r, err := os.Open("fixtures/walrus.jpg")

	if err != nil {
		t.Fatalf("Failed to open test image, %v", err)
	}

	defer r.Close()

	wr := io.Discard

	fnum, err := PrepareTag("FNumber", "11/1")

	if err != nil {
		t.Fatalf("Failed to prepare tag, %v", err)
	}

	xres, err := PrepareTag("XResolution", "72/1")

	if err != nil {
		t.Fatalf("Failed to prepare tag, %v", err)
	}

	props := map[string]interface{}{
		"FNumber":     fnum,
		"XResolution": xres,
		// "GPSLatitude": "2 13 8",
	}

	// update_test.go:52: Failed to update EXIF data, failed to set GPSLatitude tag: value not encodable: [float64] [37.61799]

	err = UpdateExif(r, wr, props)

	if err != nil {
		t.Fatalf("Failed to update EXIF data, %v", err)
	}

	// TO DO: READ AND VALIDATE TAGS
}
