package tags

import (
	"fmt"
	"testing"
)

func TestSupportedTags(t *testing.T) {

	supported, err := SupportedTags()

	if err != nil {
		t.Fatalf("Failed to determine supported tags, %v", err)
	}

	fmt.Println(len(supported))
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
