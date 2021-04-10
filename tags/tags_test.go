package tags

import (
	_ "fmt"
	"testing"
)

func TestSupportedTags(t *testing.T) {

	supported, err := SupportedTags()

	if err != nil {
		t.Fatalf("Failed to determine supported tags, %v", err)
	}

	// go run -mod vendor cmd/tags-supported/main.go | wc -l
	expected_count := 86

	supported_count := len(supported)

	if supported_count != expected_count {
		t.Fatalf("Unexpected count for supported tags. Expected %d, but got %d", expected_count, supported_count)
	}

}

func TestIsSupported(t *testing.T) {

	supported := []string{
		"Artist",
		"Make",
		"SubSecTimeOriginal",
	}

	unsupported := []string{
		"FlashpixVersion",
		"SubjectDistance",
		"ImageWidth",
	}

	for _, tag := range supported {

		ok, err := IsSupported(tag)

		if err != nil {
			t.Fatalf("Failed to determine whether '%s' is supported, %v", tag, err)
		}

		if !ok {
			t.Fatalf("Tag '%s' is not supported but should be", tag)
		}
	}

	for _, tag := range unsupported {

		ok, err := IsSupported(tag)

		if err != nil {
			t.Fatalf("Failed to determine whether '%s' is supported, %v", tag, err)
		}

		if ok {
			t.Fatalf("Tag '%s' is supported but should not be (yet)", tag)
		}
	}
}
