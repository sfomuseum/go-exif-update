package www

import (
	"embed"
)

//go:embed *.html wasm/* css/* javascript/* images/*
var FS embed.FS
