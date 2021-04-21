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
	}

	// update_test.go:52: Failed to update EXIF data, failed to set GPSLatitude tag: value not encodable: [float64] [37.61799]

	err = UpdateExif(r, wr, props)

	if err != nil {
		t.Fatalf("Failed to update EXIF data, %v", err)
	}

	// TO DO: READ AND VALIDATE TAGS
}

func TestUpdateExifGPS(t *testing.T) {

	r, err := os.Open("fixtures/walrus.jpg")

	if err != nil {
		t.Fatalf("Failed to open test image, %v", err)
	}

	defer r.Close()

	wr := io.Discard

	lat := 37.61799
	lon := -122.384864

	gps_lat, err := PrepareDecimalGPSLatitudeTag(lat)

	if err != nil {
		t.Fatalf("Failed to prepare GPSLatitudeTag, %v", err)
	}

	gps_lon, err := PrepareDecimalGPSLongitudeTag(lon)

	if err != nil {
		t.Fatalf("Failed to prepare GPSLatitudeTag, %v", err)
	}

	gps_lat_ref, err := PrepareDecimalGPSLatitudeRefTag(lat)

	if err != nil {
		t.Fatalf("Failed to prepare GPSLatitudeRefTag, %v", err)
	}

	gps_lon_ref, err := PrepareDecimalGPSLongitudeRefTag(lon)

	if err != nil {
		t.Fatalf("Failed to prepare GPSLatitudeRefTag, %v", err)
	}

	props := map[string]interface{}{
		"GPSLatitude":     gps_lat,
		"GPSLatitudeRef":  gps_lat_ref,
		"GPSLongitude":    gps_lon,
		"GPSLongitudeRef": gps_lon_ref,
	}

	err = UpdateExif(r, wr, props)

	if err != nil {
		t.Fatalf("Failed to update EXIF data, %v", err)
	}

	// TO DO: READ AND VALIDATE TAGS
}
