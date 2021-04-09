// package www provides a embedded filesystem (`embed.FS`) containing a simple web application to demonstrate the `update_exif.wasm` WebAssembly binary.
package www

import (
	"embed"
)

//go:embed *.html wasm/* css/* javascript/* images/*
var FS embed.FS
