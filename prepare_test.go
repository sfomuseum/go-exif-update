package update

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"testing"
)

func TestPrepareDecimalGPSLatitudeTag(t *testing.T) {

	ref := "N"
	lat := 37.61799

	r, err := PrepareDecimalGPSLatitudeTag(lat)

	if err != nil {
		t.Fatalf("Failed to prepare latitude '%v', %v", lat, err)
	}

	gd, err := exif.NewGpsDegreesFromRationals(ref, r.([]exifcommon.Rational))

	if err != nil {
		t.Fatalf("Failed to create NewGpsDegreesFromRationals, %v", err)
	}

	new_lat := gd.Degrees + (gd.Minutes / 60.0) + (gd.Seconds / 3600.0)

	// we expect the conversion to be inexact beyond 3 decimal points

	if fmt.Sprintf("%0.3f", lat) != fmt.Sprintf("%0.3f", new_lat) {
		t.Fatalf("Failed to convert latitude")
	}

}

func TestPrepareDecimalGPSLongitudeTag(t *testing.T) {

	ref := "W"
	lon := -122.384864

	r, err := PrepareDecimalGPSLongitudeTag(lon)

	if err != nil {
		t.Fatalf("Failed to prepare longitude '%v', %v", lon, err)
	}

	gd, err := exif.NewGpsDegreesFromRationals(ref, r.([]exifcommon.Rational))

	if err != nil {
		t.Fatalf("Failed to create NewGpsDegreesFromRationals, %v", err)
	}

	new_lon := gd.Degrees + (gd.Minutes / 60.0) + (gd.Seconds / 3600.0)

	new_lon = -new_lon // W

	// we expect the conversion to be inexact beyond 3 decimal points

	if fmt.Sprintf("%0.3f", lon) != fmt.Sprintf("%0.3f", new_lon) {
		t.Fatalf("Failed to convert longitude")
	}

}
