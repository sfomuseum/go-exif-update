package update

import (
	"testing"
)

func TestAppendGPSPropertiesWithLatitudeAndLongitude(t *testing.T) {

	lat := 37.61799
	lon := -122.384864

	props := make(map[string]interface{})

	err := AppendGPSPropertiesWithLatitudeAndLongitude(props, lat, lon)

	if err != nil {
		t.Fatalf("Failed to append GPS properties, %v", err)
	}
}
