package tags

import (
	_ "embed"
	"github.com/dsoprea/go-exif/v3/common"
	"gopkg.in/yaml.v2"
	"sync"
)

//go:embed tags_data.yaml
var tags_data []byte

var tags_init sync.Once
var tags_supported []string

var UnsupportedTypes []exifcommon.TagTypePrimitive
var UnsupportedTypesString []string

func init() {

	// These are not supported yet because I am not sure what if any value-wrangling
	// we need to do to ensure they are recorded as valid EXIF tags.
	// (20210409/thisisaaronland)

	UnsupportedTypes = []exifcommon.TagTypePrimitive{
		exifcommon.TypeRational,
		exifcommon.TypeSignedRational,
		exifcommon.TypeShort,
		exifcommon.TypeLong,
		exifcommon.TypeSignedLong,
		exifcommon.TypeUndefined,
	}

	unsupported_str := make([]string, len(UnsupportedTypes))

	for idx, t := range UnsupportedTypes {
		unsupported_str[idx] = t.String()
	}

	UnsupportedTypesString = unsupported_str
}

// https://github.com/dsoprea/go-exif/blob/de2141190595193aa097a2bf3205ba0cf76dc14b/tags.go#L189
type encodedTag struct {
	// id is signed, here, because YAML doesn't have enough information to
	// support unsigned.
	Id       int    `yaml:"id"`
	Name     string `yaml:"name"`
	TypeName string `yaml:"type_name"`
}

func IsSupported(t string) (bool, error) {

	supported, err := SupportedTags()

	if err != nil {
		return false, err
	}

	for _, this_t := range supported {

		if this_t == t {
			return true, nil
		}
	}

	return false, nil
}

func SupportedTags() ([]string, error) {

	var tags_err error

	tags_func := func() {

		tags_supported = make([]string, 0)

		encodedIfds := make(map[string][]encodedTag)

		err := yaml.Unmarshal(tags_data, encodedIfds)

		if err != nil {
			tags_err = err
			return
		}

		for _, ifdtags := range encodedIfds {

			for _, t := range ifdtags {

				for _, ts := range UnsupportedTypesString {

					if t.TypeName == ts {
						continue
					}
				}

				tags_supported = append(tags_supported, t.Name)
			}
		}
	}

	tags_init.Do(tags_func)

	if tags_err != nil {
		return nil, tags_err
	}

	return tags_supported, nil
}
