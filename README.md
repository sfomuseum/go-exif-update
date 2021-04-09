# go-exif-update

Go package to for updating EXIF data in JPEG files.

This is a thin wrapper around the dsoprea's [go-exif](https://github.com/dsoprea/go-exif) and [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure) packages and includes command-line tools for updating the EXIF data JPEG files using key-value parameters as well as a WebAssembly (wasm) binary for updating EXIF data in JavaScript (or other languages that support wasm binaries).

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-exif-update.svg)](https://pkg.go.dev/github.com/sfomuseum/go-exif-update)

## Example

```
package main

import (
	"flag"
	"github.com/sfomuseum/go-exif-update"
	"log"
	"os"
)

func main() {

	exif_props := map[string]interface{}{
		"Artist": "Bob",
		"Copyright": "SFO Museum",
	}
	
	for _, path := range paths {

		fh, _ := os.Open(path)
		defer fh.Close()

		update.UpdateExif(fh, os.Stdout, exif_props)
	}
}

```

_Error handling removed for the sake of brevity._

## Tools

```
$> make cli
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go
go build -mod vendor -o bin/tags-is-supported cmd/tags-is-supported/main.go
go build -mod vendor -o bin/tags-supported cmd/tags-supported/main.go
go build -mod vendor -o bin/update-exif cmd/update-exif/main.go
go build -mod vendor -o bin/server cmd/update-exif-server/main.go
```

As part of the build process for tools the two WebAssembly (wasm) binaries that are used by the `update-exif-server` tool are compiled and placee. You can also build the wasm binaries separately using the `wasm` Makefile target:

```
$> make wasm
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go
```

### tags-is-supported

Command-line tool for indicating whether a named EXIF tag is supported by the sfomuseum/go-exif-update package.

```
$> ./bin/tags-is-supported -h
Command-line tool for indicating whether a named EXIF tag is supported by the sfomuseum/go-exif-update package.

Usage:
	./bin/tags-is-supported tag(N) tag(N) tag(N)
```

For example:

```
```

### tags-supported

### update-exif

### update-exif-server

## See also

* https://github.com/dsoprea/go-exif
* https://github.com/dsoprea/go-jpeg-image-structure
* https://exiftool.org/TagNames/EXIF.html