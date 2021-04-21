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
}
