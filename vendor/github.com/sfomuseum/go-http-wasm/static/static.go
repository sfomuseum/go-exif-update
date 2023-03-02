// package static provides an `embed.FS` containing JavaScript used by the go-http-wasm package.
package static

import (
	"embed"
)

//go:embed javascript/*
var FS embed.FS
